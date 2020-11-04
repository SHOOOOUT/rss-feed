package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron"
)

type Slack struct {
	Text      string `json:"text"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
}

func main() {
	//ZennのrssフィードURL
	zenn := "https://zenn.dev/feed"

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(zenn)

	items := feed.Items

	//.envからフィードを流したいチャンネルのURLを取得
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}
	IncomingURL := os.Getenv("INCOMINGURL")

	c := cron.New()
	c.AddFunc("0 9 * * * *", func() {
		for i, _ := range items {
			params := Slack{
				Text:      "「" + items[i].Title + "」" + "\n" + items[i].Link + "\n" + "---------------------------------------",
				Username:  "鉄人28号",
				IconEmoji: ":28:",
			}

			jsonparams, _ := json.Marshal(params)
			resp, _ := http.PostForm(
				IncomingURL,
				url.Values{"payload": {string(jsonparams)}},
			)
			defer resp.Body.Close()
		}
	})

	c.Start()
	defer c.Stop()

	select {}
}
