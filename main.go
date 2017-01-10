package main

import (
    // "encoding/json"
    "log"
    "fmt"
    "net/http"
    "strconv"

    // mux "github.com/gorilla/mux"
    // uuid "github.com/satori/go.uuid"
    socketio "github.com/googollee/go-socket.io"
    redis "github.com/garyburd/redigo/redis"

    votingViews "./views/voting_views"
    adminViews "./views/admin_views"
    channelViews "./views/channel_views"
    viewHelpers "./views/view_helpers"
)

var pubSubChannel string = "CrowdPlay_NodeCommunication"
var Rconn *redis.Conn

func main() {
    // Setup connection to redis
    Rconn, err := redis.Dial("tcp", ":6379")
    if err != nil {
        fmt.Printf("no connection possible! try opening redis in terminal")
        return
    }

    // Give other functions a reference to what Rconn is
    votingViews.Rconn = Rconn
    adminViews.Rconn = Rconn
    channelViews.Rconn = Rconn

    // Begin serving REST API routes
    fmt.Printf("redis connection made!")
    // apirouter := mux.NewRouter()
    // log.Fatal(http.ListenAndServe(":8080", apirouter))
    // fmt.Printf("set up mux router")
    io, err := socketio.NewServer(nil)
    if err != nil {
        log.Fatal(err)
        return
    }

    // Setup socket server
    io.On("connection", func(so socketio.Socket) {
		log.Println("established a new connection")

        // ALL USER ACTION
        so.On("joinChannel", func(channel string) string {
            existence, err := channelViews.ChannelExists(channel)
            if err != nil {
                fmt.Println(err)
                return fmt.Sprintf("error: %s", err)
            }
            if !existence {
                errMsg := fmt.Sprintf("error: channel %s does not exist, create it!", channel)
                fmt.Println(errMsg)
                return errMsg
            }
            so.Join(channel)
            msg := fmt.Sprintf("joined %s", channel)
            fmt.Printf(msg)
            return fmt.Sprintf("success: %s", msg)
        })
        so.On("createChannel", func(channel string) string {
            existence, err := channelViews.ChannelExists(channel)
            if err != nil {
                fmt.Println(err)
                return fmt.Sprintf("error: %s", err)
            }
            if existence {
                errMsg := fmt.Sprintf("error: channel %s already exists, cannot be created", channel)
                fmt.Println(errMsg)
                return errMsg
            }
            so.Join(channel)
            return fmt.Sprintf("success: channel %s created", channel)
        })
        so.On("upvote", func(jsonMsg string) string {
            // unmarshall the JSON
            songData, err := viewHelpers.JsonStringUnmarshal(jsonMsg)
            if err != nil {
                fmt.Println(err)
                return fmt.Sprintf("error: %s", err)
            }
            channel := songData["channel"]
            songIDstr := songData["songID"]
            songID, err := strconv.ParseInt(songIDstr, 10, 64)
            if err != nil {
                fmt.Println(err)
                return fmt.Sprintf("error: %s", err)
            }
            msg := fmt.Sprintf("Upvote for channel: %v, songID: %v", channel, songID)
            log.Println(msg)

            // should we check the room from what we know about user, or from request?

            votingViews.UpVote(channel, songID)
            return fmt.Sprintf("sucess: %s", msg)
            // publish via redis PUBSUB channel
        })
        so.On("downvote", func(jsonMsg string) string {
            // unmarshall the JSON
            songData, err := viewHelpers.JsonStringUnmarshal(jsonMsg)
            if err != nil {
                fmt.Println(err)
                return fmt.Sprintf("error: %s", err)
            }
            channel := songData["channel"]
            songIDstr := songData["songID"]
            songID, err := strconv.ParseInt(songIDstr, 10, 64)
            if err != nil {
                fmt.Println(err)
                return fmt.Sprintf("error: %s", err)
            }
            log.Println("Downvote for channel: %v, songID: %v", channel, songID)

            votingViews.DownVote(channel, songID)
            return "success: downvote successful"
            // publish via redis PUBSUB channel
        })
        so.On("addSong", func(jsonMsg string) string {
            // unmarshall the JSON
            jsonMap, err := viewHelpers.JsonStringUnmarshal(jsonMsg)
            if err != nil {
                fmt.Println(err)
                return fmt.Sprintf("error: %s", err)
            }
            channel := jsonMap["channel"]
            votingViews.AddSong(channel, jsonMsg)
            return fmt.Sprintf("song added successfully to channel %s", channel)
            // publish via redis PUBSUB channel
        })

        // ADMIN ONLY ACTIONS
        so.On("pause", func(jsonMsg string) {
            jsonMap, err := viewHelpers.JsonStringUnmarshal(jsonMsg)
            if err != nil {
                fmt.Println(err)
                return
            }
            channel := jsonMap["channel"]
            adminViews.PauseSong(channel)
            // publish via redis PUBSUB channel
        })
        so.On("play", func(jsonMsg string) {
            jsonMap, err := viewHelpers.JsonStringUnmarshal(jsonMsg)
            if err != nil {
                fmt.Println(err)
                return
            }
            channel := jsonMap["channel"]
            adminViews.PlaySong(channel)
            // publish via redis PUBSUB channel
        })
        so.On("skipSong", func(jsonMsg string) {
            jsonMap, err := viewHelpers.JsonStringUnmarshal(jsonMsg)
            if err != nil {
                fmt.Println(err)
                return
            }
            channel := jsonMap["channel"]
            adminViews.SkipSong(channel)
            // publish via redis PUBSUB channel
        })


        so.On("disconnection", func() {
            // change metadata of the channel this was connected with to reflect change
			log.Println("on disconnect")
		})
	})
	io.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

    http.Handle("/socket.io/", io)
	http.Handle("/", http.FileServer(http.Dir("./ui")))
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
