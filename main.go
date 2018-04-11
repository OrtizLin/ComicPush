package main

import (
	"ComicNotify/crawler"
	"ComicNotify/linebot"
	"fmt"
	"net/http"
	"os"
)

func main() {
	app, err := NewLineBot(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
		os.Getenv("APP_BASE_URL"),
	)
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/start", crawler.Start)
	http.HandleFunc("/callback", app.Callback)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		fmt.Println(err)
	}
}
