package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
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

	go writer(ws)
	reader(ws)
}

func writer(ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}
