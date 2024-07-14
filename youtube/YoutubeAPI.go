package youtube

import "github.com/YoutubeVideoStats/models"

type YouTubeAPI interface {
	GetVideoByID(videoID string) (models.Video, error)
	GetChannelByID(channelID string) (models.Channel, error)
	GetVideoCommentsByID(videoID string) ([]models.Comment, error)
}
