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
	"strconv"
	"strings"
	"time"
)

const BaseAddress = "http://www.manhuagui.com"

var count uint64

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

/*Line bot*/
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
				fmt.Println(err)
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

/*Line bot*/

/*爬蟲*/

func countUpdater() {
	for {
		CrawlAndSent()
		count++
		time.Sleep(600 * time.Second)
	}
}

func CrawlAndSent() {
	fmt.Println("START CRAWL...")
	var BOOL = true
	var queryString string = ""
	pageCount := 1

	now := time.Now()
	local1, err1 := time.LoadLocation("Asia/Chongqing")
	if err1 != nil {
		fmt.Println(err1)
	}
	time_one := now.In(local1)

	var comics []NewComic
	for BOOL {
		if queryString == "" {
			queryString = BaseAddress + "/list/update.html"
		} else {
			queryString = BaseAddress + "/list/update_p" + strconv.Itoa(pageCount) + ".html"
		}
		doc, err := goquery.NewDocument(queryString)
		if err != nil {
			fmt.Println("ERROR SHOWS UP")
			log.Fatal(err)
			BOOL = false
		} else {
			doc.Find("li").Each(func(i int, s *goquery.Selection) {
				comic := NewComic{}
				title, existed := s.Find("a.bcover").Attr("title")
				if existed {
					date := s.Find("span.updateon").Text()
					result := strings.Replace(date, "更新于：", "", -1)[:10]
					fmt.Println("找到關於 " + title + " 的資料, 更新時間為 : " + result)
					if title == "约定的梦幻岛" || title == "一拳超人" || title == "进击的巨人" || title == "ONE PIECE航海王" || title == "Dr.STONE" || title == "猎人" || title == "排球少年！！" || title == "中華小廚師" {
						if result == time_one.Format("2006-01-02") {
							fmt.Println(title + "在近日內有更新！！")
							comic.Title = title
							comic.Date = result
							href, _ := s.Find("a.bcover").Attr("href")
							comic.Link = GetLink(href)
							comics = append(comics, comic)
						}
					}
					if time_one.Format("2006-01-02") != result {
						BOOL = false
					}
				}
			})
		}

		pageCount = pageCount + 1
	}

	fmt.Println("查到" + strconv.Itoa(len(comics)) + "筆資料")
	//有最新更新時, 檢查DB是否已經存在
	if strconv.Itoa(len(comics)) != "0" {
		session, errs := mgo.Dial(os.Getenv("DBURL"))
		if errs != nil {
			panic(errs)
		}
		defer session.Close()
		c := session.DB("xtest").C("commicdata")
		c2 := session.DB("xtest").C("commicuser")
		for i := 0; i < len(comics); i++ {
			result := comics[i]
			err := c.Find(bson.M{"link": comics[i].Link}).One(&result)
			if err != nil {
				//新的連載,放入DB
				c.Insert(&NewComic{comics[i].Title, comics[i].Link, comics[i].Date})
				//發送至群組
				app, err := NewLineBot(
					os.Getenv("ChannelSecret"),
					os.Getenv("ChannelAccessToken"),
					os.Getenv("APP_BASE_URL"),
				)
				if err != nil {
					fmt.Println(err)
				}
				//搜尋所有Line token 並發送訊息
				result := User{}
				iter := c2.Find(nil).Iter()
				for iter.Next(&result) {
					message := comics[i].Title + "\n" + comics[i].Link
					if _, err := app.bot.PushMessage(result.UserID, linebot.NewTextMessage(message)).Do(); err != nil {
					}
				}
			} else {
				//已經存在DB 故不在重複發送
				fmt.Println("EXIST ALREADY")
			}
		}
	}
}
func GetLink(link string) (r string) {
	doc, err := goquery.NewDocument(BaseAddress + link)
	if err != nil {
		fmt.Println(err)
	}

	doc.Find("li.status").Each(func(i int, s *goquery.Selection) {
		href, existed := s.Find("a").Attr("href")
		if existed {
			r = BaseAddress + href
		}
	})
	return r

}

/*爬蟲*/

func main() {
	app, err := NewLineBot(
		os.Getenv("ChannelSecret"),
		os.Getenv("ChannelAccessToken"),
		os.Getenv("APP_BASE_URL"),
	)
	if err != nil {

		fmt.Println(err)
	}
	go countUpdater()
	http.HandleFunc("/wakeup", WakeUp)
	http.HandleFunc("/callback", app.Callback)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {

		fmt.Println(err)
	}
}

//wake up heroku server
func WakeUp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}
