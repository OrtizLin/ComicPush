package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/line/line-bot-sdk-go/linebot"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"time"
)

const BaseAddress = "http://www.manhuagui.com"

// line_bot app
type LineBot struct {
	bot         *linebot.Client
	appBaseURL  string
	downloadDir string
}

type NewComic struct {
	Title string
	Link  string
	Date  string
}

type User struct {
	UserID string
}

// NewLineBot function
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

//wake up heroku server
func WakeUp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}
func SendMessage(message string) {
	app, err := NewLineBot(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
		os.Getenv("APP_BASE_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("commicuser")

	result := User{}
	iter := c.Find(nil).Iter()
	for iter.Next(&result) {
		if _, err := app.bot.PushMessage(result.UserID, linebot.NewTextMessage(message)).Do(); err != nil {
		}
	}

}
func FindUpdate() []NewComic {
	//today's date
	log.Println("STARTING CRAWLER ...")
	now := time.Now()
	local1, err1 := time.LoadLocation("")
	if err1 != nil {
		log.Println(err1)
	}
	local2, err2 := time.LoadLocation("Asia/Chongqing")
	if err2 != nil {
		log.Println(err2)
	}
	time_one := now.In(local1)
	time_two := now.In(local2)
	log.Println(time_one.Format("2006-01-02 15:04:05"))
	log.Println(time_two.Format("2006-01-02 15:04:05"))

	var comics []NewComic

	doc, err := goquery.NewDocument(BaseAddress + "/update")
	if err != nil {
		log.Fatal(err)
		log.Println("ERROR SHOWS UP")
	}

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		comic := NewComic{}
		title, existed := s.Find("a.cover").Attr("title")
		if existed {
			if title == "约定的梦幻岛" || title == "一拳超人" || title == "进击的巨人" || title == "ONE PIECE航海王" || title == "Dr.STONE" || title == "猎人" {
				date := s.Find("span.dt").Find("em").Text()
				if date == time_one.Format("2016-01-02") || date == time_two.Format("2016-01-02") {
					comic.Title = title
					comic.Date = date
					href, _ := s.Find("a.cover").Attr("href")
					comic.Link = GetLink(href)
					comics = append(comics, comic)
				}
			}
		}
	})
	return comics

}

func GetLink(link string) (r string) {
	doc, err := goquery.NewDocument(BaseAddress + link)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("li.status").Each(func(i int, s *goquery.Selection) {
		href, existed := s.Find("a").Attr("href")
		if existed {
			r = BaseAddress + href
		}
	})
	return r

}
func countUpdater() {
	for {
		CrawlAndSend()
		count++
		time.Sleep(600 * time.Second)
	}
}
func CrawlAndSend() {
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("commicdata")
	var comics = FindUpdate()
	for i := 0; i < len(comics); i++ {
		result := comics[i]
		err := c.Find(bson.M{"link": comics[i].Link}).One(&result)
		if err != nil {
			c.Insert(&NewComic{comics[i].Title, comics[i].Link, comics[i].Date})
			SendMessage(comics[i].Title + "\n" + comics[i].Link)
		} else {
			log.Println("EXIST ALREADY")
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
		CrawlAndSend()
	default:

		user := User{}
		user.UserID = source.UserID

		session, errs := mgo.Dial(os.Getenv("DBURL"))
		if errs != nil {
			panic(errs)
		}
		defer session.Close()
		c := session.DB("xtest").C("commicuser")

		//check if userId exist.
		err := c.Find(bson.M{"userid": user.UserID}).One(&user)
		if err != nil {
			errs = c.Insert(&User{user.UserID})
			if errs != nil {
				log.Fatal(errs)
			} else {
				if _, err := app.bot.ReplyMessage(
					replyToken,
					linebot.NewTextMessage("恭喜您已訂閱連載報報。\n當有最新連載發行時將會第一時間通知您！"),
				).Do(); err != nil {
					return err
				}
			}
		} else {
			if _, err := app.bot.ReplyMessage(
				replyToken,
				linebot.NewTextMessage("您已經訂閱囉！"),
			).Do(); err != nil {
				return err
			}
		}

	}
	return nil

}

func main() {
	app, err := NewLineBot(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
		os.Getenv("APP_BASE_URL"),
	)
	if err != nil {
		log.Fatal(err)
	}
	go countUpdater()

	http.HandleFunc("/wakeup", WakeUp)
	http.HandleFunc("/callback", app.Callback)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
