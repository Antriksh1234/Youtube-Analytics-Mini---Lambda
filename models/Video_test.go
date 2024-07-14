package models

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/youtube/v3"
	"testing"
)

func TestParseVideoListResponse(t *testing.T) {
	tests := []struct {
		Name     string
		Input    youtube.VideoListResponse
		Expected Video
	}{
		{
			Name: "Empty response List",
			Input: youtube.VideoListResponse{
				Items: []*youtube.Video{},
			},
			Expected: Video{},
		},
		{
			Name: "Non empty Response List",
			Input: youtube.VideoListResponse{
				Items: []*youtube.Video{
					{
						Statistics: &youtube.VideoStatistics{
							CommentCount: 1000,
							ViewCount:    1000000,
							LikeCount:    1000,
						},
						Snippet: &youtube.VideoSnippet{
							Title:       "My first Video",
							ChannelId:   "mockID",
							Description: "mock Video Description",
							Thumbnails: &youtube.ThumbnailDetails{
								High: &youtube.Thumbnail{
									Url: "mockThumbnailURL_HIGH",
								},
								Standard: &youtube.Thumbnail{
									Url: "mockThumbnailURL_STANDARD",
								},
								Medium: &youtube.Thumbnail{
									Url: "mockThumbnailURL_MEDIUM",
								},
							},
						},
					},
				},
			},
			Expected: Video{
				Title:            "My first Video",
				ChannelID:        "mockID",
				Description:      "mock Video Description",
				Likes:            1000,
				Views:            1000000,
				NumberOfComments: 1000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			input := tt.Input
			video := ParseVideoListResponse(&input)
			assert.Equal(t, tt.Expected.Title, video.Title)
			assert.Equal(t, tt.Expected.Description, video.Description)
			assert.Equal(t, tt.Expected.Views, video.Views)
			assert.Equal(t, tt.Expected.Likes, video.Likes)
			assert.Equal(t, tt.Expected.NumberOfComments, video.NumberOfComments)
			assert.Equal(t, tt.Expected.ChannelID, video.ChannelID)
		})
	}
}
