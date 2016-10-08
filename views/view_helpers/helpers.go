package viewHelpers

import (
    "fmt"
    "encoding/json"
)

func GenChannelKey(channelName string) string {
    return fmt.Sprintf("Channel::%s_Meta", channelName)
}

func GenChannelSongQueueKey(channelName string) string {
    return fmt.Sprintf("Channel::%s_SongQueue", channelName)
}

func GenSongKeys(songID int64) (string, string) {
    songMetaKey := fmt.Sprintf("Song::%d_Meta", songID)
    songVoteKey := fmt.Sprintf("Song::%d_VoteCount", songID)
    return songMetaKey, songVoteKey
}

func GenNowPlayingKey(channelName string) string {
    return fmt.Sprintf("channel::%s_nowplaying")
}

func JsonStringUnmarshal(jsonstr string) (map[string]string, error) {
    data := make(map[string]string)
    err := json.Unmarshal([]byte(jsonstr), &data)
    if err != nil {
        fmt.Println(err)
        return data, err
    }
    return data, nil
}
