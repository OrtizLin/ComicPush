package crawler

import (
	"ComicNotify/db"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"time"
)

const BaseAddress = "http://www.manhuagui.com"

var count uint64

type NewComic struct {
	Title string
	Link  string
	Date  string
}

func Start() {
	go countUpdater()
}
func countUpdater() {
	for {
		crawl()
		count++
		time.Sleep(600 * time.Second)
	}
}

//Comic crawler
func crawl() {
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
	//有最新更新時, 檢查DB是否已經存在
	if strconv.Itoa(len(comics)) != "0" {
		db.CheckAlreadySent(comics)
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
