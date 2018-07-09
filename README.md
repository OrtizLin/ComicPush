# ComicPush

Web crawler - Saving time for search serial comic update every week.

# LINE Bot Push API

Chose developer trial can send push message unlimited。

![螢幕快照 2018-03-22 下午5.28.30](https://i.imgur.com/1AXmw88.png)

# How to use it？

Add this bot as friend , and send any message , server will get your client ID and save to database.

Once comic update , bot will send you a link.

![螢幕快照 2018-03-22 下午5.33.38](https://i.imgur.com/UFtngZR.png)

# Others

```
if title == "约定的梦幻岛" || title == "一拳超人" || title == "进击的巨人" 
|| title == "ONE PIECE航海王" || title == "Dr.STONE" || title == "猎人" || title == "排球少年！！" {
		if result == time_one.Format("2006-01-02") {
					fmt.Println(title + "在近日內有更新！！")
					comic.Title = title
					comic.Date = result
					href, _ := s.Find("a.bcover").Attr("href")
					comic.Link = GetLink(href)
					comics = append(comics, comic)
				}
		}
```

Comics: 約定的夢幻島 , 一拳超人 , 進擊的巨人 , 海賊王 , Dr.Stone , 獵人 , 排球少年, 七原罪, 亞人, 食戟之靈.
