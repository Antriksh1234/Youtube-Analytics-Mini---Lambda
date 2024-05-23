package test

import (
	"log"
	"sync"
	"time"

	"github.com/YoutubeVideoStats/types"
)

type SentimentAPIMock struct{}

func (s *SentimentAPIMock) GetSentiment(text string) string {
	time.Sleep(time.Millisecond * 100)
	return "POSITIVE"
}

func GetSentimentOfComment(s *SentimentAPIMock, comment string, sentimentChannel chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	sentimentChannel <- s.GetSentiment(comment)
}

func (s *SentimentAPIMock) GetSentimentOfComments(video types.Video) (string, map[string]int) {
	start := time.Now()

	log.Println("GetSentimentOfComments: START")
	SentimentStats := make(map[string]int)

	sentimentChannel := make(chan string)

	wg := &sync.WaitGroup{}
	for _, comment := range video.Comments {
		wg.Add(1)
		go GetSentimentOfComment(s, comment.CommentText, sentimentChannel, wg)
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
	max := 0
	maxSentiment := "POSITVE"

	for sentiment, countOfAppearance := range sentimentMap {
		if countOfAppearance > max || (countOfAppearance == max && maxSentiment != "POSITIVE") {
			maxSentiment = sentiment
			max = countOfAppearance
		}
	}

	return maxSentiment
}
