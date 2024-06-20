package types

type Comment struct {
	CommentText string `json:"commentText"`
	Likes       int64  `json:"likes"`
	Commenter   string `json:"commenter"`
}
