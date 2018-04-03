package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

// RegisterRoutes registers routes
func RegisterRoutes(routes *mux.Router) {
	routes.HandleFunc("/messages", getMessagesHandler).Methods("GET")
	routes.HandleFunc("/messages/{id}", getMessageHandler).Methods("GET")
}

func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	session := GetDatabaseConnection()
	defer session.Close()
	collection := GetDatabaseCollection(session, "messages")

	var result []HTTPMessage
	collection.Find(bson.M{}).Limit(10).All(&result)

	content, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, string(content))
}

func getMessageHandler(w http.ResponseWriter, r *http.Request) {
	session := GetDatabaseConnection()
	defer session.Close()
	collection := GetDatabaseCollection(session, "messages")

	vars := mux.Vars(r)
	requestID := vars["id"]

	var result HTTPMessage
	collection.Find(bson.M{"_id": requestID}).One(&result)

	content, err := json.Marshal(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, string(content))
}
