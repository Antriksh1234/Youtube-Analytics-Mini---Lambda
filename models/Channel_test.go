package models

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/youtube/v3"
	"testing"
)

func TestParseChannelListResponse(t *testing.T) {
	tests := []struct {
		Name     string
		Input    youtube.ChannelListResponse
		Expected Channel
	}{
		{
			Name: "Empty response List",
			Input: youtube.ChannelListResponse{
				Items: []*youtube.Channel{},
			},
			Expected: Channel{},
		},
		{
			Name: "Non empty Response List",
			Input: youtube.ChannelListResponse{
				Items: []*youtube.Channel{
					{
						Statistics: &youtube.ChannelStatistics{
							SubscriberCount: 143000,
						},
						Snippet: &youtube.ChannelSnippet{
							Title:       "Tech with Lucy",
							Description: "mock Channel Description",
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
			Expected: Channel{
				ChannelName:      "Tech with Lucy",
				ChannelSubs:      143000,
				ChannelThumbnail: "mockThumbnailURL_HIGH",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			input := tt.Input
			channel := ParseChannelListResponse(&input)
			assert.Equal(t, tt.Expected.ChannelThumbnail, channel.ChannelThumbnail)
			assert.Equal(t, tt.Expected.ChannelSubs, channel.ChannelSubs)
			assert.Equal(t, tt.Expected.ChannelName, channel.ChannelName)
		})
	}
}
