package sentiment

import (
	"encoding/json"
	"fmt"
	"github.com/YoutubeVideoStats/models"
	"log"
	"os"
	"time"
)

type SentimentAPIMock struct{}

func (s *SentimentAPIMock) GetSentiment(text string) string {
	data, err := os.ReadFile("../sentiment/testdata/sentimentMap.json")
	if err != nil {
		fmt.Println(err)
		return "NEUTRAL"
	}

	var sentimentMap map[string]string
	err = json.Unmarshal(data, &sentimentMap)
	if err != nil {
		return "NEUTRAL"
	}
	if val, ok := sentimentMap[text]; ok {
		return val
	}

	return "NEUTRAL"
}

func (s *SentimentAPIMock) GetSentimentOfComments(comments []models.Comment) (string, map[string]int) {
	start := time.Now()

	log.Println("GetSentimentOfComments: START")
	SentimentStats := make(map[string]int)

	for _, comment := range comments {
		sentiment := s.GetSentiment(comment.CommentText)
		if val, ok := SentimentStats[sentiment]; ok {
			SentimentStats[sentiment] = val + 1
		} else {
			SentimentStats[sentiment] = 1
		}
	}

	Sentiment := getDominantSentimentOfComments(SentimentStats)
	log.Println("GetSentimentOfComments: END; Time Taken: ", time.Since(start))

	return Sentiment, SentimentStats
}
