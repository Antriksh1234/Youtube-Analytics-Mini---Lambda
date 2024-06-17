package api

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/YoutubeVideoStats/types"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Youtube struct {
	Service *youtube.Service
}

var APIKey = os.Getenv("API_KEY")

func InitializeYoutubeService() (YoutubeAPI, error) {
	service, err := youtube.NewService(context.TODO(), option.WithAPIKey(APIKey))
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	Youtube := Youtube{
		Service: service,
	}

	return Youtube, err
}

func (y Youtube) GetVideoByID(videoID string) types.Video {
	start := time.Now()
	log.Println("GetVideoByID: START")
	video := types.Video{}

	call := y.Service.Videos.List([]string{"statistics", "snippet"}).Id(videoID)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making API call to retrieve video details: %v", err)
		return video
	}

	if len(response.Items) == 0 {
		log.Println("Video not found")
		return video
	}

	statistics := response.Items[0].Statistics
	snippet := response.Items[0].Snippet

	video.Title = snippet.Title
	video.ThumbnailURL = snippet.Thumbnails.High.Url
	video.Views = statistics.ViewCount
	video.Likes = statistics.LikeCount
	video.ChannelID = snippet.ChannelId
	video.Description = snippet.Description
	video.NumberOfComments = statistics.CommentCount

	log.Println("GetVideoByID: END; Time taken: ", time.Since(start))
	return video
}

func (y Youtube) GetChannelByID(channelID string) types.Channel {
	start := time.Now()

	log.Println("GetChannelByID: START")
	youtubeChannel := types.Channel{}

	call := y.Service.Channels.List([]string{"snippet", "statistics"}).Id(channelID)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error fetching channel details: %v", err)
		return youtubeChannel
	}

	if len(response.Items) == 0 {
		log.Println("Channel not found")
		return youtubeChannel
	}

	snippet := response.Items[0].Snippet
	statistics := response.Items[0].Statistics

	youtubeChannel.ChannelName = snippet.Title
	youtubeChannel.ChannelThumbnail = snippet.Thumbnails.High.Url
	youtubeChannel.ChannelSubs = statistics.SubscriberCount

	log.Println("GetChannelByID: END; Time Taken: ", time.Since(start))
	return youtubeChannel
}

func (y Youtube) GetVideoCommentsByID(videoID string) ([]types.Comment, error) {
	start := time.Now()

	log.Println("GetVideoCommentsByID: START")
	var Comments []types.Comment
	nextPageToken := ""

	for {
		call := y.Service.CommentThreads.List([]string{"snippet"}).VideoId(videoID).PageToken(nextPageToken)
		response, err := call.Do()
		if err != nil {
			log.Printf("Error fetching comments for video %s: %v\n", videoID, err)
			return nil, err
		}

		for _, item := range response.Items {
			commentText := item.Snippet.TopLevelComment.Snippet.TextDisplay
			commentAuthor := item.Snippet.TopLevelComment.Snippet.AuthorDisplayName
			commentLikes := item.Snippet.TopLevelComment.Snippet.LikeCount

			Comment := types.Comment{
				CommentText: commentText,
				Commenter:   commentAuthor,
				Likes:       commentLikes,
			}
			Comments = append(Comments, Comment)

			if len(Comments) > 200 {
				return Comments, nil
			}
		}

		if response.NextPageToken == "" {
			break // No more pages, exit loop
		}

		nextPageToken = response.NextPageToken
	}

	log.Println("GetVideCommentsByID: END; Time Taken: ", time.Since(start))
	return Comments, nil
}
