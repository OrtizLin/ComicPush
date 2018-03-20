package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
)

const BaseAddress = "http://www.manhuagui.com"

type NewComic struct {
	Title string
	Link  string
	Date  string
}

func FindUpdate() []NewComic {

	var comics []NewComic
	doc, err := goquery.NewDocument(BaseAddress + "/update")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		comic := NewComic{}
		title, existed := s.Find("a.cover").Attr("title")
		if existed {
			if title == "约定的梦幻岛" || title == "一拳超人" || title == "进击的巨人" || title == "ONE PIECE航海王" || title == "Dr.STONE" {
				comic.Title = title
				comic.Date = s.Find("span.dt").Find("em").Text()
				href, _ := s.Find("a.cover").Attr("href")
				comic.Link = GetLink(href)
				comics = append(comics, comic)
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
func main() {
	log.Println(FindUpdate())
}
