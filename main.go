package main

import (
	"ComicNotify/bot"
	"ComicNotify/crawler"
	"fmt"
	"net/http"
	"os"
)

/*爬蟲*/

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
	http.HandleFunc("/callback", bot.Callback)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		fmt.Println(err)
	}
}
