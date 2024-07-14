package main

import (
	"context"
	"github.com/YoutubeVideoStats/app"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
)

func HandleRequest(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Hello youtube Gateway event! %s\n", event.Body)
	newApp, err := app.NewApp()
	if err != nil {
		return newApp.CreateErrorResponse(http.StatusInternalServerError), err
	}
	return newApp.FetchVideoFeedback(event)
}

func main() {
	lambda.Start(HandleRequest)
}
