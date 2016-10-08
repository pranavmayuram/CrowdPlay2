package votingViews

import (
    viewHelpers "../view_helpers"
    redis "github.com/garyburd/redigo/redis"
    "encoding/json"
    "fmt"
    "strconv"
    "errors"
)

var Rconn redis.Conn

// Bool indicates if song was added. If return false,
func AddSong(channelName string, songMetaDataString string) (bool, error) {
    // decode the JSON from string to a JSON object
    songData := make(map[string]string)
    err := json.Unmarshal([]byte(songMetaDataString), &songData)
    if err != nil {
        fmt.Println(err)
        return false, err
    }
    songIDstr := songData["songID"]
    songID, err := strconv.ParseInt(songIDstr, 10, 64)
    exists, err := SongExists(channelName, songID)
    if err != nil {
        fmt.Println(err)
        return false, err
    } else if !exists {
        errmsg := fmt.Sprintf("Song %v already exists, upvoting instead of adding.", songID)
        err := errors.New(errmsg)
        fmt.Println(errmsg)
        UpVote(channelName, songID)
        return false, err
    }
    // Consider turning the JSON object into class (could also be a point of validation)
    channelQueueKey := viewHelpers.GenChannelSongQueueKey(channelName)
    songMetaKey, songVoteKey := viewHelpers.GenSongKeys(songID)
    Rconn.Send("HSET", channelQueueKey, songMetaKey, songMetaDataString)
    voteCreation, err := redis.Bool(Rconn.Do("HSET", channelQueueKey, songVoteKey, 1))
    if err != nil {
        fmt.Println(err)
        return false, err
    }
    return voteCreation, nil
}



func UpVote(channelName string, songID int64) (int, error) {
    // do an increment via Redis
    return AdjustVote(channelName, songID, 1)
}

func DownVote(channelName string, songID int64) (int, error) {
    // do a decrement via Redis
    return AdjustVote(channelName, songID, -1)
}



func AdjustVote(channelName string, songID int64, voteWeight int) (int, error) {
    channelQueueKey := viewHelpers.GenChannelSongQueueKey(channelName)
    _, songVoteKey := viewHelpers.GenSongKeys(songID)
    newVoteCount, err := redis.Int(Rconn.Do("HINCRBY", channelQueueKey, songVoteKey, voteWeight))
    if err != nil {
        fmt.Println(err)
        return -1, err
    }
    return newVoteCount, nil
}

func SongExists(channelName string, songID int64) (bool, error){
    channelQueueKey := viewHelpers.GenChannelSongQueueKey(channelName)
    songMetaKey, _ := viewHelpers.GenSongKeys(songID)
    existence, err := redis.Bool(Rconn.Do("HEXISTS", channelQueueKey, songMetaKey))
    if err != nil {
        fmt.Println(err)
        return false, err
    }
    return existence, nil
}
