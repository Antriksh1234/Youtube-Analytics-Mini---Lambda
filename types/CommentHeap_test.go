package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetTopCommentsByLikes(t *testing.T) {

	//Case when the number of comments is less than 100
	NumberOfComments := 20
	Comments := makeDummyComments(NumberOfComments)

	TopComments := GetTopCommentsByLikes(Comments)
	assert.Equal(t, len(TopComments), NumberOfComments)
	assert.Equal(t, int64(NumberOfComments-1), TopComments[0].Likes)

	//Case when the number of comments is greater than 100
	NumberOfComments = 120
	Comments = makeDummyComments(NumberOfComments)
	TopComments = GetTopCommentsByLikes(Comments)

	assert.Equal(t, 100, len(TopComments)) //We want to pick only top 100 by likes
	assert.Equal(t, int64(NumberOfComments-1), TopComments[0].Likes)
	assert.Equal(t, int64(20), TopComments[len(TopComments)-1].Likes)

	//Case when there are no comments
	NumberOfComments = 0
	Comments = makeDummyComments(NumberOfComments)
	TopComments = GetTopCommentsByLikes(Comments)

	assert.Equal(t, 0, len(TopComments))
}

func makeDummyComments(NumberOfComments int) []Comment {
	Comments := make([]Comment, NumberOfComments)
	for i := 0; i < NumberOfComments; i++ {
		Comments[i] = Comment{
			CommentText: "",
			Likes:       int64(i),
		}
	}

	return Comments
}
