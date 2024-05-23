package test

import (
	"testing"

	"github.com/YoutubeVideoStats/app"
)

func Test_App(t *testing.T) {
	app := app.App{}

	app.YoutubeAPI = &YoutubeAPIMock{}
	app.SentimentAPI = &SentimentAPIMock{}

}
