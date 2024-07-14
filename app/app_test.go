package app

import (
	"fmt"
	"github.com/YoutubeVideoStats/models"
	"github.com/YoutubeVideoStats/sentiment"
	"github.com/YoutubeVideoStats/youtube"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func Test_extractVideoURL(t *testing.T) {
	tests := []struct {
		Name     string
		Input    string
		Expected string
	}{
		{
			Name:     "Valid URL 1",
			Input:    "https://www.youtube.com/watch?v=xhgHeAhxizE",
			Expected: "xhgHeAhxizE",
		},
		{
			Name:     "Valid URL 2",
			Input:    "https://www.youtube.com/watch?v=mockID",
			Expected: "mockID",
		},
		{
			Name:     "Valid URL 3",
			Input:    "https://youtu.be/2dc63Lhdr-I?si=j5MV7B_12D5q",
			Expected: "2dc63Lhdr-I",
		},
		{
			Name:     "Invalid URL 4",
			Input:    "Invalid",
			Expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			videoID, _ := extractVideoID(tt.Input)
			assert.Equal(t, tt.Expected, videoID)
		})
	}
}

func TestGetAPIs(t *testing.T) {
	app := MockApp{
		YoutubeAPI:   &youtube.YoutubeAPIMock{},
		SentimentAPI: &sentiment.SentimentAPIMock{},
	}

	assert.NotEqual(t, nil, app.GetYoutubeAPI())
	assert.NotEqual(t, nil, app.GetSentimentAPI())
}

func TestCreateErrorResponse(t *testing.T) {
	app := App{}

	resp := app.CreateErrorResponse(http.StatusInternalServerError)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	resp2 := app.CreateErrorResponse(http.StatusBadRequest)
	assert.Equal(t, http.StatusBadRequest, resp2.StatusCode)

}

func Test_getYoutubeVideoFromYoutubeAPI(t *testing.T) {
	app := MockApp{
		YoutubeAPI:   &youtube.YoutubeAPIMock{},
		SentimentAPI: &sentiment.SentimentAPIMock{},
	}

	videoInfo, err := getVideoInfoFromYouTubeAPI(&app, "iAX4n0ZlTUY")
	assert.Equal(t, nil, err)
	fmt.Println(videoInfo)

	assert.Equal(t, 4, len(videoInfo.Comments))
}

func Test_getChannelInfoFromYoutubeAPI(t *testing.T) {
	app := MockApp{
		YoutubeAPI:   &youtube.YoutubeAPIMock{},
		SentimentAPI: &sentiment.SentimentAPIMock{},
	}

	channelInfo, err := app.GetYoutubeAPI().GetChannelByID("123")

	assert.Equal(t, nil, err)
	assert.Equal(t, "Tech With Lucy", channelInfo.ChannelName)
}

func Test_fetchFeedbackOfVideo(t *testing.T) {
	mockApp := MockApp{
		YoutubeAPI:   &youtube.YoutubeAPIMock{},
		SentimentAPI: &sentiment.SentimentAPIMock{},
	}

	tests := []struct {
		Name             string
		Request          events.APIGatewayProxyRequest
		ExpectedFeedback models.Feedback
		ExpectedError    bool
	}{
		{
			Name: "Empty Body",
			Request: events.APIGatewayProxyRequest{
				Body: "",
			},
			ExpectedFeedback: models.Feedback{},
			ExpectedError:    true,
		},
		{
			Name: "Empty VideoID",
			Request: events.APIGatewayProxyRequest{
				Body: "www.youtube.com",
			},
			ExpectedFeedback: models.Feedback{},
			ExpectedError:    true,
		},
		{
			Name: "Video ID Available",
			Request: events.APIGatewayProxyRequest{
				Body: "https://youtu.be/iAX4n0ZlTUY?si=0QYFFAI7mu8YXJvB",
			},
			ExpectedFeedback: models.Feedback{
				Video: models.Video{
					Title:            "Build With Me: PartyRock AI Application | AWS Project",
					ThumbnailURL:     "url",
					Views:            1230,
					Likes:            123,
					NumberOfComments: 4,
					ChannelID:        "123",
					Description:      "Video description",
				},
				Channel: models.Channel{
					ChannelName:      "Tech With Lucy",
					ChannelSubs:      143000,
					ChannelThumbnail: "channel_high_thumbnail",
				},
			},
			ExpectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			resp, err := fetchVideoFeedback(&mockApp, tt.Request)
			assert.Equal(t, tt.ExpectedError, err != nil)
			fmt.Println(resp)
		})
	}
}
