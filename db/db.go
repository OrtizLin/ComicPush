package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
)

type NewComic struct {
	Title string
	Link  string
	Date  string
}

type User struct {
	UserID string
}

func GetAllUser() (user []User) {
	var users []User
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c2 := session.DB("xtest").C("commicuser")
	result := User{}
	iter := c2.Find(nil).Iter()
	for iter.Next(&result) {
		users = append(users, result)
	}
	return users
}

func CheckComicInDB(title, link, date string) bool {
	result := NewComic{}
	session, errs := mgo.Dial(os.Getenv("DBRUL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("commicdata")
	err := c.Find(bson.M{"link": link}).One(&result)
	if err != nil {
		c.Insert(&NewComic{title, link, date})
		return false
	} else {
		return true
	}
}

func CheckRegistered(userID string) bool {
	user := User{}
	user.UserID = userID

	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("commicuser")

	//check if userId exist.
	err := c.Find(bson.M{"userid": user.UserID}).One(&user)
	if err != nil {
		// add to database and send message.
		errs = c.Insert(&User{user.UserID})
		if errs != nil {
			log.Print(err)
		}
		return false
	} else {
		// already in database.
		return true
	}
}

func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}