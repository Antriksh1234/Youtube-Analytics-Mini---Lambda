package youtube

import (
	"encoding/json"
	"fmt"
	"github.com/YoutubeVideoStats/models"
	"os"
)

type YoutubeAPIMock struct{}

var responseMap map[string]string

func (y *YoutubeAPIMock) GetVideoByID(videoID string) (models.Video, error) {
	responseMap = map[string]string{
		"iAX4n0ZlTUY": "../youtube/testdata/video/videoResponse1.json",
	}
	if path, ok := responseMap[videoID]; ok {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}
		var video models.Video
		_ = json.Unmarshal(data, &video)
		return video, nil
	} else {
		return models.Video{}, nil
	}
}

func (y *YoutubeAPIMock) GetChannelByID(channelID string) (models.Channel, error) {
	responseMap = map[string]string{
		"123": "../youtube/testdata/channel/channelResponse1.json",
	}
	if path, ok := responseMap[channelID]; ok {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}
		var channel models.Channel
		_ = json.Unmarshal(data, &channel)
		return channel, nil
	} else {
		return models.Channel{}, nil
	}
}

func (y *YoutubeAPIMock) GetVideoCommentsByID(videoID string) ([]models.Comment, error) {
	responseMap = map[string]string{
		"iAX4n0ZlTUY": "../youtube/testdata/comment/commentResponse1.json",
	}
	if path, ok := responseMap[videoID]; ok {
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}
		var channel []models.Comment
		_ = json.Unmarshal(data, &channel)
		return channel, nil
	} else {
		return []models.Comment{}, nil
	}
}
