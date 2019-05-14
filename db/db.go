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

type RegisteredComic struct {
	ComicName string
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
	log.Print(title, date, link)
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("commicdata")
	err := c.Find(bson.M{"link": link}).One(&result)
	if err != nil {
		log.Print("插入資料庫!")
		log.Print(title, date, link)
		c.Insert(&NewComic{title, link, date})
		return false
	} else {
		return true
	}
}

func PrintRegistered() string {
	var results []RegisteredComic
	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("registercomic")
	err := c.Find(nil).All(&results)
 	if err != nil {
		panic(errs)
	}
 	str := ""
	for i := 0; i < len(results); i++ {
		str += results[i].ComicName + " ,"
	}
	return str

}

func CheckRegisteredComic(comicName string) bool {
	comic := RegisteredComic{}
	comic.ComicName = comicName

	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("registercomic")
	err := c.Find(bson.M{"comicname": comic.ComicName}).One(&comic)

	if err != nil {
		return false // comic is not exist
	} else { 
		return true // comic is exist
	}	
}

func RegisterComic(comicName string) {
	if CheckRegisteredComic(comicName) == false {
	comic := RegisteredComic{}
	comic.ComicName = comicName

	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("registercomic")

	errs = c.Insert(&RegisteredComic{comic.ComicName})
	if errs != nil{
		log.Print(errs)
	}
}
}

func CheckRegisteredUser(userID string) bool {
	user := User{}
	user.UserID = userID

	session, errs := mgo.Dial(os.Getenv("DBURL"))
	if errs != nil {
		panic(errs)
	}
	defer session.Close()
	c := session.DB("xtest").C("commicuser")

	//check if userId is exist.
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
