package crawler

import (
	"ComicNotify/bot"
	"ComicNotify/db"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type NewComic struct {
	Title string
	Link  string
	Date  string
}

const BaseAddress = "http://www.manhuagui.com"

var count uint64

func Start(w http.ResponseWriter, r *http.Request) {
	go countUpdater()
}

func countUpdater() {
	for {
		crawlAndSent()
		count++
		time.Sleep(600 * time.Second)
	}
}

func crawlAndSent() {
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
					if title == "约定的梦幻岛" || title == "一拳超人" || title == "进击的巨人" || title == "ONE PIECE航海王" || title == "Dr.STONE" || title == "猎人" || title == "排球少年！！" || title == "中華小廚師" || title == "妖神记" {
						if result == time_one.Format("2006-01-02") {
							fmt.Println(title + "在近日內有更新！！")
							comic.Title = title
							comic.Date = result
							href, _ := s.Find("a.bcover").Attr("href")
							comic.Link = getLink(href)
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
		for i := 0; i < len(comics); i++ {
			if db.CheckComicInDB(comics[i].Title, comics[i].Link, comics[i].Date) {
				//已經存在DB 故不在重複發送
			} else {
				str := comics[i].Title + "\n" + comics[i].Link
				bot.PushMessage(str)
			}
		}
	}
}
func getLink(link string) (r string) {
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
