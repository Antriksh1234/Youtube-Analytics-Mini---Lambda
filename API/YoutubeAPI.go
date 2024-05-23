package api

import "github.com/YoutubeVideoStats/types"

type YoutubeAPI interface {
	GetVideoByID(videoID string) types.Video
	GetChannelByID(channelID string) types.Channel
	GetVideoCommentsByID(videoID string) ([]types.Comment, error)
}
