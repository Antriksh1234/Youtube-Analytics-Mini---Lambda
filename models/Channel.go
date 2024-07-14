package models

import (
	"google.golang.org/api/youtube/v3"
	"log"
)

type Channel struct {
	ChannelName      string `json:"channel"`
	ChannelSubs      uint64 `json:"channelSubs"`
	ChannelThumbnail string `json:"channelThumbnail"`
}

func ParseChannelListResponse(response *youtube.ChannelListResponse) Channel {
	youtubeChannel := Channel{}

	if len(response.Items) == 0 {
		log.Println("Channel not found")
		return youtubeChannel
	}

	snippet := response.Items[0].Snippet
	statistics := response.Items[0].Statistics

	youtubeChannel.ChannelName = snippet.Title
	youtubeChannel.ChannelThumbnail = snippet.Thumbnails.High.Url
	youtubeChannel.ChannelSubs = statistics.SubscriberCount

	return youtubeChannel
}
