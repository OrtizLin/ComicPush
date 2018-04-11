package db

import (
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
func CheckAlreadySent(comics NewComic) bool {
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("commicdata")
	result := NewComic{}
	err := c.Find(bson.M{"link": comics.Link}).One(&result)
	if err != nil {
		//新的連載,放入DB
		c.Insert(&NewComic{comics.Title, comics.Link, comics.Date})
		if err != nil {
			fmt.Println(err)
		}
		return false
	} else {
		//已經存在DB
		return true
	}

}

//Search all users
func SearchUsers()a []User {
c2 := session.DB("xtest").C("commicuser")
	//搜尋所有Line token 並發送訊息
var users []User
	userone := User{}
	iter := c2.Find(nil).Iter()
	for iter.Next(&userone) {
		users = append(users,,userone)
	}
	return users
}
