package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/YoutubeVideoStats/models"
	"github.com/YoutubeVideoStats/sentiment"
	yt "github.com/YoutubeVideoStats/youtube"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type App struct {
	YoutubeAPI   yt.YouTubeAPI
	SentimentAPI sentiment.SentimentAPI
}

func (app *App) CreateErrorResponse(statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body: fmt.Sprintf("%+v", struct {
			Message string
		}{
			"Lambda faced error while fetching results",
		}),
	}
}

func NewApp() (IApp, error) {
	mApp := App{}
	err := mApp.initializeSentimentAPI()
	if err != nil {
		return &mApp, err
	}

	err = mApp.initializeYoutubeAPI()
	if err != nil {
		return &mApp, err
	}

	return &mApp, nil
}

func (app *App) initializeYoutubeAPI() error {
	var APIKey = os.Getenv("API_KEY")
	service, err := youtube.NewService(context.TODO(), option.WithAPIKey(APIKey))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	Youtube := yt.Youtube{
		Service: service,
	}

	app.YoutubeAPI = Youtube
	return err
}

func (app *App) initializeSentimentAPI() error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		log.Fatalf("Could not initialize the Comprehend Client %v", err)
		return err
	}

	ComprehendClient := comprehend.NewFromConfig(cfg)
	AWSComprehend := sentiment.AWSComprehened{
		ComprehendClient: ComprehendClient,
	}

	app.SentimentAPI = AWSComprehend
	return err
}

func (app *App) FetchVideoFeedback(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return fetchVideoFeedback(app, event)
}

func fetchVideoFeedback(app IApp, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	videoFeedback := models.Feedback{} //This will be the YouTube response

	//Get the Video URL from the event
	videoURL := event.Body

	//Let's extract the Video ID to pass it to our YouTube API
	videoID, err := extractVideoID(videoURL)

	if err != nil {
		log.Println("Could not get Video ID from Payload, ", err.Error())
		return app.CreateErrorResponse(http.StatusBadRequest), err
	}

	video, err := getVideoInfoFromYouTubeAPI(app, videoID)
	if err != nil {
		log.Println("Could not fetch Video for Provided Video ID ", err.Error())
		return app.CreateErrorResponse(http.StatusInternalServerError), err
	}

	channel, err := getChannelInfoFromYouTubeAPI(app, video.ChannelID)

	if err != nil {
		log.Println("Could not fetch Channel for Provided channel ID ", err.Error())
		return app.CreateErrorResponse(http.StatusInternalServerError), err
	}

	videoFeedback.Channel = channel
	videoFeedback.Video = video

	//Let's marshal this feedback and send as a response
	data, err := json.Marshal(videoFeedback)
	if err != nil {
		log.Println("Could not marshal youtube video feedback: ", err.Error())
		return app.CreateErrorResponse(http.StatusInternalServerError), err
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

func getChannelInfoFromYouTubeAPI(app IApp, id string) (models.Channel, error) {
	return app.GetYoutubeAPI().GetChannelByID(id)
}

func getVideoInfoFromYouTubeAPI(app IApp, videoID string) (models.Video, error) {
	videoChan := make(chan models.Video)
	commentsChan := make(chan []models.Comment)

	go getVideoInfo(app, videoID, videoChan)
	go getVideoComments(app, videoID, commentsChan)

	video := <-videoChan
	comments := <-commentsChan

	video.Comments = comments

	video.Sentiment, video.SentimentStats = app.GetSentimentAPI().GetSentimentOfComments(comments)
	return video, nil
}

func getVideoComments(app IApp, videoID string, commentsChan chan<- []models.Comment) {
	comments, _ := app.GetYoutubeAPI().GetVideoCommentsByID(videoID)
	commentsChan <- comments
}

func getVideoInfo(app IApp, videoID string, videoChan chan<- models.Video) {
	videoInfo, _ := app.GetYoutubeAPI().GetVideoByID(videoID)
	videoChan <- videoInfo
}

func (app *App) GetSentimentAPI() sentiment.SentimentAPI {
	return app.SentimentAPI
}

func (app *App) GetYoutubeAPI() yt.YouTubeAPI {
	return app.YoutubeAPI
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
