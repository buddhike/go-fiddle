package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go-fiddle/internal/kafkaclient"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	kafkaClient := kafkaclient.NewConsumer(func(msg *kafka.Message) {
		log.Printf("Message received %v\n%s\n", msg.TopicPartition, string(msg.Value))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

	kafkaClient.Close()
}
