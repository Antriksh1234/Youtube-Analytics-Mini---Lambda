package models

import (
	"google.golang.org/api/youtube/v3"
	"log"
)

type Video struct {
	Title            string    `json:"title"`
	ThumbnailURL     string    `json:"thumbnailURL"`
	Views            uint64    `json:"views"`
	Likes            uint64    `json:"likes"`
	NumberOfComments uint64    `json:"numberOfComments"`
	ChannelID        string    `json:"channelID"`
	Comments         []Comment `json:"comments"`
	Description      string    `json:"description"`

	Sentiment      string         `json:"sentiment"`
	SentimentStats map[string]int `json:"sentimentStats"`
}

func ParseVideoListResponse(response *youtube.VideoListResponse) Video {
	video := Video{}
	if len(response.Items) == 0 {
		log.Println("Video not found")
		return video
	}

	statistics := response.Items[0].Statistics
	snippet := response.Items[0].Snippet

	video.Title = snippet.Title
	video.ThumbnailURL = snippet.Thumbnails.High.Url
	video.Views = statistics.ViewCount
	video.Likes = statistics.LikeCount
	video.ChannelID = snippet.ChannelId
	video.Description = snippet.Description
	video.NumberOfComments = statistics.CommentCount

	return video
}
