package linebot

import (
	"ComicNotify/db"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
)

type User struct {
	UserID string
}

type LineBotStruct struct {
	bot         *linebot.Client
	appBaseURL  string
	downloadDir string
}

func NewLineBot(channelSecret, channelToken, appBaseURL string) (*LineBotStruct, error) {
	bots, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, err
	}
	return &LineBotStruct{
		bot:         bots,
		appBaseURL:  appBaseURL,
		downloadDir: "testing",
	}, nil
}

func (app *LineBotStruct) Callback(w http.ResponseWriter, r *http.Request) {
	events, err := app.bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err := app.handleText(message, event.ReplyToken, event.Source); err != nil {
					log.Print(err)
				}
			default:
				log.Printf("Unknown message: %v", message)
			}
		default:
			log.Printf("Unknown event: %v", event)
		}
	}
}

func (app *LineBotStruct) handleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	switch message.Text {

	default:
		if db.RegisterUser(source.UserID) {
			//already exist
			if _, err := app.bot.ReplyMessage(
				replyToken,
				linebot.NewTextMessage("您已經訂閱囉！"),
			).Do(); err != nil {
				return err
			}
		} else {
			if _, err := app.bot.ReplyMessage(
				replyToken,
				linebot.NewTextMessage("恭喜您已訂閱連載報報。\n當有最新連載發行時將會第一時間通知您！"),
			).Do(); err != nil {
				return err
			}
		}

	}
	return nil

}

func PushMessage(userID, message string) {
	//發送至群組
	app, err := NewLineBot(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
		os.Getenv("APP_BASE_URL"),
	)
	if _, err := app.bot.PushMessage(userID, linebot.NewTextMessage(message)).Do(); err != nil {
	}
}
