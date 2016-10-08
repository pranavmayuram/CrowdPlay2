package models

// General song MetaData that should be included for any song
type SongMetaData struct {
    SongID int64
    SongName string
    Genre string
    ImageLocation string
    LengthInSeconds int
}

// For songs that are waiting to be played in the queue
type SongEnqueued struct {
    SongInfo *SongMetaData
    VoteCount int
}

// For songs that are waiting to be played in the queue
type SongNowPlaying struct {
    SongInfo *SongMetaData
    Paused int
    Progress int
}

// Holds metadata about the entire channel
type Channel struct {
    ChannelName string
    SongsPlayed int
    NumContributors int
}

func NewChannelModel(channelName string) *Channel {
    return &Channel{ChannelName: channelName, SongsPlayed: 0, NumContributors: 0}
}

type VoteMessage struct {
    UIViewChanged string
    SongID int64
    NewVoteCount int
}

type SongEnqueuedMessage struct {
    UIViewChanged string
    SongInfo *SongEnqueued
}

type ChannelMetadataChanged struct {

}
