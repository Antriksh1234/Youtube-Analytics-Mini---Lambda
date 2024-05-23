package awsservices

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/YoutubeVideoStats/api"
	YoutubeVideoStatsTypes "github.com/YoutubeVideoStats/types"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	"github.com/aws/aws-sdk-go-v2/service/comprehend/types"
)

type AWSComprehened struct {
	ComprehendClient *comprehend.Client
}

func InitializeSentimentAPI() (api.SentimentAPI, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		log.Fatalf("Could not initialize the Comprehend Client %v", err)
	}

	ComprehendClient := comprehend.NewFromConfig(cfg)
	AWSComprehend := AWSComprehened{
		ComprehendClient: ComprehendClient,
	}

	return AWSComprehend, err
}

func (awscomprehend AWSComprehened) GetSentiment(text string) string {
	comprehendClient := awscomprehend.ComprehendClient
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

func GetSentimentOfComment(awscomprehend AWSComprehened, comment string, sentimentChannel chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	sentimentChannel <- awscomprehend.GetSentiment(comment)
}

func (awscomprehend AWSComprehened) GetSentimentOfComments(video YoutubeVideoStatsTypes.Video) (string, map[string]int) {
	start := time.Now()

	log.Println("GetSentimentOfComments: START")
	SentimentStats := make(map[string]int)

	sentimentChannel := make(chan string)

	wg := &sync.WaitGroup{}
	for _, comment := range video.Comments {
		wg.Add(1)
		go GetSentimentOfComment(awscomprehend, comment.CommentText, sentimentChannel, wg)
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
