package bot

import (
	"comic-push-robot/db"
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
				app.UserRegister("", event.ReplyToken, event.Source)
			}
		default:
			app.UserRegister("", event.ReplyToken, event.Source)
		}
	}
}

func (app *LineBot) handleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {
	switch message.Text {
	case "查詢":
		str := db.PrintRegistered()
		replymessage := "目前加入的漫畫有 : " + str
		if _, err := app.bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage(replymessage),
		).Do(); err != nil {
			return err
		}
	default:
		app.UserRegister(message.Text, replyToken, source)
	}
	return nil
}

func (app *LineBot) UserRegister(message string, replyToken string, source *linebot.EventSource) error {

	//str := ""
	userID := source.UserID
	
	if userID == "" {
		userID = source.GroupID
	}

	if userID == "" {
		userID = source.RoomID
	}
	
	if db.CheckRegisteredUser(userID) {
		//str = "看看此作者其他的機器人吧！"
	} else {
		//str = "恭喜您已訂閱連載報報。\n當有最新連載發行時將會第一時間通知您！"
	}

	if source.UserID == os.Getenv("MASTER_UUID") && message != "" {
		if message == "test" { // Test giant image nessage.


		users := db.GetAllUser()
		for i := 0; i < len(users); i++ {
			if _, err := app.bot.PushMessage(users[i].UserID, linebot.NewImagemapMessage(
					"https://i.imgur.com/ITTfpLC.png",
					"Try this giant image",
					linebot.ImagemapBaseSize{Width: 1040, Height: 2080},
					linebot.NewURIImagemapAction("https://www.manhuagui.com/", linebot.ImagemapArea{X: 0, Y: 0, Width: 1040, Height: 2080}),
				),
			).Do(); err != nil {
			}
		}
		} else {
		if db.CheckRegisteredComic(message) {
			//str = message + "已在你的資料庫裡"
		} else {
			db.RegisterComic(message)
			//str = "將 " + message + " 加入資料庫！"
		}
	}
	}
	//log.Print(str)
	// if str == "看看此作者其他的機器人吧！" {
	// 	var columns []*linebot.CarouselColumn

	// 	column1 := linebot.NewCarouselColumn(
	// 		"https://i.imgur.com/l3Cdj6B.png", "表特爆報", "只要 PTT 表特版上出現爆文，立即通知使用者，讓你不再錯過精彩文章。",
	// 		linebot.NewURITemplateAction("立即點我追蹤", "https://line.me/R/ti/p/%40qss3676v"),
	// 	)
	// 	column2 := linebot.NewCarouselColumn(
	// 		"https://i.imgur.com/q7Jq9WK.png", "空汙爆報", "客製化訂閱你家或工作場所附近的空汙觀測站，於每日8點及下午6點進行空汙推播通知。",
	// 		linebot.NewURITemplateAction("立即點我追蹤", "https://line.me/R/ti/p/%40zhm1865k"),
	// 	)
	// 	columns = append(columns, column1)
	// 	columns = append(columns, column2)

	// 	template := linebot.NewCarouselTemplate(columns...)
	// 	if _, err := app.bot.ReplyMessage(
	// 		replyToken,
	// 		linebot.NewTemplateMessage(str, template),
	// 	).Do(); err != nil {
	// 		return err

	// 	}
	// } else {
		// if _, err := app.bot.ReplyMessage(
		// 	replyToken,
		// 	linebot.NewTextMessage(str),
		// ).Do(); err != nil {
		// 	return err
		// }
	// }

	//IMPLEMENT NEW FEATURE
	//KEEP IMPLEMENT NEW FEATURE
	return nil
}
