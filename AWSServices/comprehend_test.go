package awsservices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getDominantSentimentOfComments(t *testing.T) {
	sentimentMap := make(map[string]int)

	//Case 1: POSITIVE is clearly the winner
	sentimentMap["POSITIVE"] = 3
	sentimentMap["NEGATIVE"] = 2
	sentimentMap["NEUTRAL"] = 2
	sentimentMap["MIXED"] = 2
	assert.Equal(t, "POSITIVE", getDominantSentimentOfComments(sentimentMap))

	//Case 2: NEGATIVE is clearly the winner
	sentimentMap["POSITIVE"] = 2
	sentimentMap["NEGATIVE"] = 3
	sentimentMap["NEUTRAL"] = 2
	sentimentMap["MIXED"] = 2
	assert.Equal(t, "NEGATIVE", getDominantSentimentOfComments(sentimentMap))

	//Case 3: NEUTRAL is clearly the winner
	sentimentMap["POSITIVE"] = 2
	sentimentMap["NEGATIVE"] = 2
	sentimentMap["NEUTRAL"] = 3
	sentimentMap["MIXED"] = 2
	assert.Equal(t, "NEUTRAL", getDominantSentimentOfComments(sentimentMap))

	//Case 4: MIXED is clearly the winner
	sentimentMap["POSITIVE"] = 2
	sentimentMap["NEGATIVE"] = 2
	sentimentMap["NEUTRAL"] = 3
	sentimentMap["MIXED"] = 4
	assert.Equal(t, "MIXED", getDominantSentimentOfComments(sentimentMap))

	//Case 5: POSITIVE ties with something else but we give it a POSITVE rating
	sentimentMap["POSITIVE"] = 4
	sentimentMap["NEGATIVE"] = 2
	sentimentMap["NEUTRAL"] = 3
	sentimentMap["MIXED"] = 4
	assert.Equal(t, "POSITIVE", getDominantSentimentOfComments(sentimentMap))
}
