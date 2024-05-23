package types

type Video struct {
	Title            string    `json:"title"`
	ThumbnailURL     string    `json:"thumbnailURL"`
	Views            uint64    `json:"views"`
	Likes            uint64    `json:"likes"`
	NumberOfComments uint64    `json:"numberOfComments"`
	ChannelID        string    `json:"channelID"`
	Comments         []Comment `json:"comments"`
	Description      string    `json:"description"`

	Sentiment      string         `json:"sentiment"`
	SentimentStats map[string]int `json:"sentimentStats"`
}
