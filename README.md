# ComicPush

為了解決每週要去搜尋最新漫畫連載的問題 , 寫一個爬蟲爬常看的線上漫畫網頁 , 並做資料的比對 , 若喜歡的漫畫更新 , 利用 Line bot push api 通知。

# Line bot push api

選擇 Developer Trial 可以使用 Line bot push message 且無推播訊息數量限制 , 一個帳戶最高可新增 50 個好友。

![螢幕快照 2018-03-22 下午5.28.30](https://i.imgur.com/1AXmw88.png)

# 如何使用？

加為好友後 , 第一次使用者只要傳送訊息 , 即可取得 User ID 並存入資料庫,有漫畫更新時, 將根據 ID 個別發送漫畫連結。

![螢幕快照 2018-03-22 下午5.33.38](https://i.imgur.com/UFtngZR.png)

# 其他

```
if title == "约定的梦幻岛" || title == "一拳超人" || title == "进击的巨人" || title == "ONE PIECE航海王" || title == "Dr.STONE" || title == "猎人" || title == "排球少年！！" {
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

將標題為 約定的夢幻島, 一拳超人, 進擊的巨人, 海賊王, Dr.Stone, 獵人, 排球少年

的漫畫列出, 若更新日期為今天, 則推送通知。