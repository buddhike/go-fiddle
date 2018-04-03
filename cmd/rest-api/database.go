package main

import (
	"go-fiddle/internal/config"

	mgo "gopkg.in/mgo.v2"
)

// GetDatabaseConnection gets mongo database connection
func GetDatabaseConnection() *mgo.Session {
	session, err := mgo.Dial(config.Get("MONGODB", "mongodb://localhost"))
	if err != nil {
		panic(err)
	}

	// session.SetMode(mgo.Monotonic, true)

	return session
}

// GetDatabaseCollection gets mongo database collection
func GetDatabaseCollection(session *mgo.Session, collectionName string) *mgo.Collection {
	return session.DB("gofiddle").C(collectionName)
}
