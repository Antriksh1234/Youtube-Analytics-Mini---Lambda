package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/YoutubeVideoStats/api"
	"github.com/YoutubeVideoStats/types"
	"github.com/aws/aws-lambda-go/events"
)

type Feedback struct {
	types.Video   //Info of the video
	types.Channel //Channel that posted the video
}

type App struct {
	YoutubeAPI   api.YoutubeAPI
	SentimentAPI api.SentimentAPI
}

func (app *App) GetFeedbackOfYoutubeVideo(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	videoFeedback := Feedback{} //This will be the API response

	//Get the Video URL from the event
	videoURL := event.Body

	//Let's extract the Video ID to pass it to our API
	videoID, err := extractVideoID(videoURL)

	if err != nil {
		log.Fatal("Could not get Video ID from Payload, ", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	if err != nil {
		log.Fatal("Could not Initialize our Youtube API, ", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	//We will do three things
	//1. Get the Video Stats from YoutubeAPI
	//2. Get the Channel stats from YoutubeAPI
	//3. Get the sentiment of comments from SentimentAPI

	youtubeVideo := app.YoutubeAPI.GetVideoByID(videoID)
	youtubeChannel := app.YoutubeAPI.GetChannelByID(youtubeVideo.ChannelID)
	youtubeVideo.Comments, _ = app.YoutubeAPI.GetVideoCommentsByID(videoID)
	youtubeVideo.NumberOfComments = uint64(len(youtubeVideo.Comments))

	//Lets just return top comments
	youtubeVideo.Comments = types.GetTopCommentsByLikes(youtubeVideo.Comments)

	youtubeVideo.Sentiment, youtubeVideo.SentimentStats = app.SentimentAPI.GetSentimentOfComments(youtubeVideo)

	videoFeedback.Channel = youtubeChannel
	videoFeedback.Video = youtubeVideo

	//Lets marshal this feedback and send as a response
	data, err := json.Marshal(videoFeedback)
	if err != nil {
		log.Fatal("Could not marshal youtube video feedback: ", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
		}, nil
	}

	//We parsed the feedback, lets send this
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"content-type": "application/json",
		},
		StatusCode: 200,
		Body:       string(data),
	}, nil
}

func extractVideoID(videoURL string) (string, error) {
	u, err := url.Parse(videoURL)
	if err != nil {
		return "", err
	}

	if u.Host == "youtu.be" {
		// Extract video ID from short-form URL
		videoID := strings.TrimPrefix(u.Path, "/")
		if videoID == "" {
			return "", fmt.Errorf("unable to extract video ID from URL")
		}
		return videoID, nil
	}

	queryParams, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}

	videoID := queryParams.Get("v")
	if videoID == "" {
		return "", fmt.Errorf("unable to extract video ID from URL")
	}

	return videoID, nil
}
