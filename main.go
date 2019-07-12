package main

import (
	"comic-push-robot/bot"
	"comic-push-robot/crawler"
	"fmt"
	"net/http"
	"os"
)

func main() {
	app, err := bot.NewLineBot(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
		os.Getenv("APP_BASE_URL"),
	)
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/wakeup", crawler.Start)
	http.HandleFunc("/callback", app.Callback)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		fmt.Println(err)
	}
}

// I Need To fix this bug.