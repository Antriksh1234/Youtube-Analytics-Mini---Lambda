package test

import "github.com/YoutubeVideoStats/types"

type YoutubeAPIMock struct{}

func (y *YoutubeAPIMock) GetVideoByID(videoID string) types.Video {
	video := types.Video{
		Title:            "Favourite Video",
		ThumbnailURL:     "Favourite Video Thumbnail",
		Views:            136000,
		Likes:            1234,
		NumberOfComments: 1342,
		ChannelID:        "123",
	}

	size := 100
	for size > 0 {
		video.Comments = append(video.Comments, types.Comment{
			CommentText: "ABCD",
			Likes:       12,
			Commenter:   "Antriksh",
		})
		size--
	}

	return video
}

func (y *YoutubeAPIMock) GetChannelByID(channelID string) types.Channel {
	return types.Channel{
		ChannelName:      "Tech With Lucy",
		ChannelSubs:      148000,
		ChannelThumbnail: "thumbnail",
	}
}

func (y *YoutubeAPIMock) GetVideoCommentsByID(videoID string) ([]types.Comment, error) {
	size := 100
	Comments := []types.Comment{}
	for size > 0 {
		Comments = append(Comments, types.Comment{
			CommentText: "ABCD",
			Likes:       12,
			Commenter:   "Antriksh",
		})
		size--
	}

	return Comments, nil
}
