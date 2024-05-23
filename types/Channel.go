package types

type Channel struct {
	ChannelName      string `json:"channel"`
	ChannelSubs      uint64 `json:"channelSubs"`
	ChannelThumbnail string `json:"channelThumbnail"`
}
