package sentiment

import "github.com/YoutubeVideoStats/models"

type SentimentAPI interface {
	GetSentiment(text string) string
	GetSentimentOfComments(v []models.Comment) (string, map[string]int)
}
