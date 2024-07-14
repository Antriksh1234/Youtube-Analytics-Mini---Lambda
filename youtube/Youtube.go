package youtube

import (
	"log"
	"time"

	"github.com/YoutubeVideoStats/models"
	"google.golang.org/api/youtube/v3"
)

type Youtube struct {
	Service *youtube.Service
}

func (y Youtube) GetVideoByID(videoID string) (models.Video, error) {
	start := time.Now()
	log.Println("getVideoByID: START")
	video := models.Video{}

	call := y.Service.Videos.List([]string{"statistics", "snippet"}).Id(videoID)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making youtube call to retrieve video details: %v", err)
		return video, err
	}

	log.Println("getVideoByID: END; Time taken: ", time.Since(start))
	return models.ParseVideoListResponse(response), nil
}

func (y Youtube) GetChannelByID(channelID string) (models.Channel, error) {
	start := time.Now()

	log.Println("getChannelByID: START")
	youtubeChannel := models.Channel{}

	call := y.Service.Channels.List([]string{"snippet", "statistics"}).Id(channelID)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error fetching channel details: %v", err)
		return youtubeChannel, err
	}

	log.Println("getChannelByID: END; Time Taken: ", time.Since(start))
	return models.ParseChannelListResponse(response), nil
}

func (y Youtube) GetVideoCommentsByID(videoID string) ([]models.Comment, error) {
	start := time.Now()

	log.Println("getVideoCommentsByID: START")
	var Comments []models.Comment
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

			Comment := models.Comment{
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
