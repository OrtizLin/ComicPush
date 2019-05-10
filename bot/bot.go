package bot

import (
	"ComicNotify/db"
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
)

type LineBot struct {
	bot         *linebot.Client
	appBaseURL  string
	downloadDir string
}

func NewLineBot(channelSecret, channelToken, appBaseURL string) (*LineBot, error) {
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, err
	}

	return &LineBot{
		bot:         bot,
		appBaseURL:  appBaseURL,
		downloadDir: "test",
	}, nil
}

func PushMessage(msg string) {

	app, err := NewLineBot(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
		os.Getenv("APP_BASE_URL"),
	)
	if err != nil {
		fmt.Println(err)
	}
	users := db.GetAllUser()
	for i := 0; i < len(users); i++ {
		if _, err := app.bot.PushMessage(users[i].UserID, linebot.NewTextMessage(msg)).Do(); err != nil {
		}
	}
}

func (app *LineBot) Callback(w http.ResponseWriter, r *http.Request) {
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

func (app *LineBot) handleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	switch message.Text {
	case "tonygrr":
		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("嫩！"),
		).Do(); err != nil {
			return err
		}
	default:
		str := ""
		if db.CheckRegisteredUser(source.UserID) {
			str = "您已經訂閱囉！" + source.UserID
		} else {
			str = "恭喜您已訂閱連載報報。\n當有最新連載發行時將會第一時間通知您！"
		}

		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(str),
		).Do(); err != nil {
			return err
		}

		if source.UserID == "U91942915be32583fc6d583cccd3c1dc8" {
			if _, errs := app.bot.ReplyMessage(
				replyToken,
				linebot.NewTextMessage("you are the master"),
			).Do(); errs != nil {
				return errs
			}
		}

	}
	return nil
}
