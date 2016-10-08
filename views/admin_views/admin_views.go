package adminViews

import (
    viewHelpers "../view_helpers"
    redis "github.com/garyburd/redigo/redis"
    // "log"
    // "encoding/json"
    "fmt"
)

var Rconn redis.Conn

func SkipSong(channelName string) (bool, error) {
    // Find the highest voted song in this channel
    // Overwrite the Now Playing doc for this channel with the new song's info
    // Potentially change the metadata doc of the channel (reflect how many songs have been played)
    return true, nil
}



func PauseSong(channelName string) (bool, error) {
    // Change status of the Now Playing doc to say paused is true
    return ChangePauseStatus(channelName, true)
}

func PlaySong(channelName string) (bool, error) {
    // Change status of the Now Playing doc to say paused is true
    return ChangePauseStatus(channelName, false)
}


func ChangePauseStatus(channelName string, pausedVal bool) (bool, error) {
    // Change status of the Now Playing doc to say paused is false
    nowPlayingKey := viewHelpers.GenNowPlayingKey(channelName)
    pauseResult, err := redis.Bool(Rconn.Do("HSET", nowPlayingKey, "paused", pausedVal))
    if err != nil {
        fmt.Println(err)
        return false, err
    }
    return pauseResult, nil
}
