package sentiment

import (
	YoutubeVideoStatsTypes "github.com/YoutubeVideoStats/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getDominantSentimentOfComments(t *testing.T) {
	tests := []struct {
		Name     string
		Input    map[string]int
		Expected string
	}{
		{
			Name: "POSITIVE is clearly the winner",
			Input: map[string]int{
				"POSITIVE": 3,
				"NEGATIVE": 2,
				"NEUTRAL":  2,
				"MIXED":    2,
			},
			Expected: "POSITIVE",
		},
		{
			Name: "NEGATIVE is clearly the winner",
			Input: map[string]int{
				"POSITIVE": 2,
				"NEGATIVE": 3,
				"NEUTRAL":  2,
				"MIXED":    2,
			},
			Expected: "NEGATIVE",
		},
		{
			Name: "NEUTRAL is clearly the winner",
			Input: map[string]int{
				"POSITIVE": 2,
				"NEGATIVE": 2,
				"NEUTRAL":  3,
				"MIXED":    2,
			},
			Expected: "NEUTRAL",
		},
		{
			Name: "MIXED is clearly the winner",
			Input: map[string]int{
				"POSITIVE": 2,
				"NEGATIVE": 2,
				"NEUTRAL":  2,
				"MIXED":    3,
			},
			Expected: "MIXED",
		},
		{
			Name: "POSITIVE ties with something else but we give it a POSITIVE rating",
			Input: map[string]int{
				"POSITIVE": 4,
				"NEGATIVE": 2,
				"NEUTRAL":  2,
				"MIXED":    4,
			},
			Expected: "POSITIVE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			assert.Equal(t, tt.Expected, getDominantSentimentOfComments(tt.Input))
		})
	}
}

func Test_fetchVideoSentiment(t *testing.T) {
	comments := []YoutubeVideoStatsTypes.Comment{
		{
			CommentText: "Good Video",
		},
		{
			CommentText: "Good Video",
		},
		{
			CommentText: "Good Video",
		},
		{
			CommentText: "Bad Video",
		},
	}

	sentiment, sentimentMap := fetchVideoSentiments(&SentimentAPIMock{}, comments)
	assert.Equal(t, "POSITIVE", sentiment)
	assert.Equal(t, 3, sentimentMap["POSITIVE"])
}
