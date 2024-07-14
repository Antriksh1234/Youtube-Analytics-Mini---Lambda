package app

import (
	"github.com/YoutubeVideoStats/sentiment"
	yt "github.com/YoutubeVideoStats/youtube"
	"github.com/aws/aws-lambda-go/events"
)

type MockApp struct {
	YoutubeAPI   yt.YouTubeAPI
	SentimentAPI sentiment.SentimentAPI
}

func (m *MockApp) CreateErrorResponse(statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
	}
}

func (m *MockApp) GetSentimentAPI() sentiment.SentimentAPI {
	return m.SentimentAPI
}

func (m *MockApp) GetYoutubeAPI() yt.YouTubeAPI {
	return m.YoutubeAPI
}

func (m *MockApp) FetchVideoFeedback(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, nil
}
