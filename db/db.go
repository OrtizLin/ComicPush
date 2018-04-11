package db

import (
	"ComicNotify/linebot"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	UserID string
}

type NewComic struct {
	Title string
	Link  string
	Date  string
}

//User register
func RegisterUser(userID string) bool {
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("commicuser")

	user := User{}
	user.UserID = source.UserID

	//check if userId exist.
	err := c.Find(bson.M{"userid": user.UserID}).One(&user)
	if err != nil {
		errs = c.Insert(&User{user.UserID})
		return false
	} else {
		return true
	}
}

//Check already sent
func CheckAlreadySent(comics []NewComic) {
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

			if err != nil {
				fmt.Println(err)
			}
			//搜尋所有Line token 並發送訊息
			results := User{}
			iter := c2.Find(nil).Iter()
			for iter.Next(&results) {
				message := comics[i].Title + "\n" + comics[i].Link
				linebot.PushMessage(results.UserID, message)
			}
		} else {
			//已經存在DB 故不在重複發送
		}
	}
}
