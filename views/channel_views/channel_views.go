package channelViews

import (
    viewHelpers "../view_helpers"
    redis "github.com/garyburd/redigo/redis"
    "encoding/json"
    "fmt"
    models "../../models"
)

var Rconn redis.Conn

func CreateChannel(channelName string) (bool, error) {
    newChan := models.NewChannelModel(channelName)
    jsonNewChan, err := json.Marshal(newChan)
    if err != nil {
        fmt.Println(err)
        return false, err
    }
    channelKey := viewHelpers.GenChannelKey(channelName)
    success, err := redis.Bool(Rconn.Do("SET", channelKey, jsonNewChan))
    if err != nil {
        fmt.Println(err)
        return false, err
    }
    return success, nil
}

func ChannelExists(channelName string) (bool, error) {
    channelKey := viewHelpers.GenChannelKey(channelName)
    existence, err := redis.Bool(Rconn.Do("EXISTS", channelKey))
    if err != nil {
        fmt.Println(err)
        return false, err
    }
    return existence, nil
}
