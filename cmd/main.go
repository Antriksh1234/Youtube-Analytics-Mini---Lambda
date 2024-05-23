package main

import (
	"context"
	"log"
	"net/http"

	awsservices "github.com/YoutubeVideoStats/AWSServices"
	"github.com/YoutubeVideoStats/api"
	"github.com/YoutubeVideoStats/app"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Hello API Gateway event! %s\n", event.Body)
	app := app.App{}

	//Initialize the Youtube API
	YoutubeAPI, err1 := api.InitializeYoutubeService()
	if err1 != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err1
	}

	//Initialize the Sentiment API
	SentimentAPI, err2 := awsservices.InitializeSentimentAPI()
	if err2 != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err2
	}

	log.Println("Initialized the Youtube API and Sentiment API Successfully")
	app.YoutubeAPI = YoutubeAPI
	app.SentimentAPI = SentimentAPI
	return app.GetFeedbackOfYoutubeVideo(event)
}

func main() {
	lambda.Start(HandleRequest)
}
