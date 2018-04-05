package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

// RegisterRoutes registers routes
func RegisterRoutes(routes *mux.Router) {
	routes.HandleFunc("/messages", getMessagesHandler).Methods("GET")
	routes.HandleFunc("/messages/{id}", getMessageHandler).Methods("GET")
	routes.HandleFunc("/ws", serveWs)
}

func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	session := GetDatabaseConnection()
	defer session.Close()
	collection := GetDatabaseCollection(session, "messages")

	var result []HTTPMessage
	collection.Find(bson.M{}).Limit(100).All(&result)

	if result == nil {
		result = []HTTPMessage{}
	}

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

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	listeners[ws] = func(msg *kafka.Message) {
		ws.WriteJSON(struct {
			Type *string `json:"type"`
		}{
			msg.TopicPartition.Topic,
		})
	}
}
