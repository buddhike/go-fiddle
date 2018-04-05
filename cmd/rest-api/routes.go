package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

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

	summary := getMessageSummary(result)

	content, err := json.Marshal(summary)

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

	listeners[ws] = func(message *HTTPMessage) {
		ws.WriteJSON(summariseMessage(*message))
	}
}

func getMessageSummary(messages []HTTPMessage) []HTTPMessageSummary {
	summary := make([]HTTPMessageSummary, len(messages))
	for i, message := range messages {
		summary[i] = summariseMessage(message)
	}
	return summary
}

func summariseMessage(message HTTPMessage) (summary HTTPMessageSummary) {
	summary = HTTPMessageSummary{
		message.ID,
		message.Request.Method,
		message.Request.URI,
		0,
	}

	if message.Response != nil {
		summary.StatusCode = message.Response.StatusCode
	}

	if !strings.HasPrefix(strings.ToLower(summary.URI), "http:") {
		for _, header := range *message.Request.Headers {
			if strings.ToLower(header.Name) == "host" {
				summary.URI = fmt.Sprintf("https://%s%s", header.Value, summary.URI)
				break
			}
		}
	}

	return
}
