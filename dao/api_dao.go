package dao

import (
	"gopkg.in/mgo.v2"
	"log"
)


const (
	DatabaseHost = "localhost:27017"
	DatabaseName = "welcome_robot"
)

func ConnectDatabase() *mgo.Database {
	session, err := mgo.Dial(DatabaseHost)
	if err != nil {
		log.Fatal(err.Error())
	}
	if session.Ping() != nil {
		log.Fatal(err)
	}
	return session.DB(DatabaseName)
}
