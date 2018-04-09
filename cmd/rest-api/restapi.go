package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"go-fiddle/cmd/config"
	"go-fiddle/cmd/internal/kafkaclient"
	"go-fiddle/cmd/internal/kafkaserver"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	listeners = make(map[*websocket.Conn]func(message *HTTPMessage))
)

func main() {
	router := mux.NewRouter()

	session := GetDatabaseConnection()
	defer session.Close()
	// session.SetMode(mgo.Monotonic, true)
	collection := GetDatabaseCollection(session, "messages")

	kafkaProducer := kafkaserver.NewProducer()

	kafkaClient := kafkaclient.NewConsumer(func(msg *kafka.Message) {
		message := string(msg.Value)

		var httpMessage *HTTPMessage
		if *msg.TopicPartition.Topic == "request" {
			requestID, request := unmarshalHTTPRequest(msg.Value)
			httpMessage = &HTTPMessage{requestID, request, nil}

			err := collection.Insert(httpMessage)

			if err != nil {
				log.Fatal(err)
			}
		} else if *msg.TopicPartition.Topic == "response" {
			requestID, response := unmarshalHTTPResponse(msg.Value)

			err := collection.FindId(requestID).One(&httpMessage)

			if err != nil {
				log.Fatal(err)
			}

			httpMessage.Response = response
			err = collection.Update(bson.M{"_id": requestID}, httpMessage)

			if err != nil {
				log.Fatal(err)
			}
		}
		log.Printf("Message received %v\n%s\n", msg.TopicPartition, message)

		if httpMessage != nil {
			for _, callback := range listeners {
				callback(httpMessage)
			}

			topic := fmt.Sprintf("%ssummary", *msg.TopicPartition.Topic)
			if jsonMessage, err := json.Marshal(summariseMessage(*httpMessage)); err == nil {
				kafkaProducer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Value:          jsonMessage,
				}, nil)
			}
		}
	})

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	port := config.Get("PORT", "8888")
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), handlers.CORS(originsOk, headersOk, methodsOk)(router)))

	kafkaClient.Close()
}
