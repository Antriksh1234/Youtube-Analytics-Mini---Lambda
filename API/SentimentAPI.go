package api

import "github.com/YoutubeVideoStats/types"

type SentimentAPI interface {
	GetSentiment(text string) string
	GetSentimentOfComments(video types.Video)(string, map[string]int)
}
