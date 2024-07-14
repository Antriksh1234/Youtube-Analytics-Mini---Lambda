package app

import (
	"github.com/YoutubeVideoStats/sentiment"
	"github.com/YoutubeVideoStats/youtube"
	"github.com/aws/aws-lambda-go/events"
)

type IApp interface {
	CreateErrorResponse(statusCode int) events.APIGatewayProxyResponse
	GetSentimentAPI() sentiment.SentimentAPI
	GetYoutubeAPI() youtube.YouTubeAPI
	FetchVideoFeedback(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}
