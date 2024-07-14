package sentiment

import (
	"context"

	YoutubeVideoStatsTypes "github.com/YoutubeVideoStats/models"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
	"log"
	"sync"
	"time"
)

type AWSComprehened struct {
	ComprehendClient *comprehend.Client
}

func (c AWSComprehened) GetSentiment(text string) string {
	comprehendClient := c.ComprehendClient
	result, err := comprehendClient.DetectSentiment(context.TODO(), &comprehend.DetectSentimentInput{
		Text:         &text,
		LanguageCode: types.LanguageCodeEn,
	})

	if err != nil {
		log.Printf("Error while fetching sentiment: %s", err.Error())
		return ""
	}

	return string(result.Sentiment)
}

func getSentimentOfComment(sentimentAPI SentimentAPI, comment string, sentimentChannel chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	sentimentChannel <- sentimentAPI.GetSentiment(comment)
}

func (c AWSComprehened) GetSentimentOfComments(comments []YoutubeVideoStatsTypes.Comment) (string, map[string]int) {
	return fetchVideoSentiments(c, comments)
}

func fetchVideoSentiments(c SentimentAPI, comments []YoutubeVideoStatsTypes.Comment) (string, map[string]int) {
	start := time.Now()

	log.Println("GetSentimentOfComments: START")
	SentimentStats := make(map[string]int)

	sentimentChannel := make(chan string)

	wg := &sync.WaitGroup{}
	for _, comment := range comments {
		wg.Add(1)
		go getSentimentOfComment(c, comment.CommentText, sentimentChannel, wg)
	}

	go func() {
		wg.Wait()
		close(sentimentChannel)
	}()

	for sentiment := range sentimentChannel {
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

func getDominantSentimentOfComments(sentimentMap map[string]int) string {
	maxAppearance := 0
	maxSentiment := "POSITVE"

	for sentiment, countOfAppearance := range sentimentMap {
		if countOfAppearance > maxAppearance || (countOfAppearance == maxAppearance && maxSentiment != "POSITIVE") {
			maxSentiment = sentiment
			maxAppearance = countOfAppearance
		}
	}

	return maxSentiment
}
