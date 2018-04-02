package kafkaserver

import (
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// NewProducer returns a new kafka producer
func NewProducer() *kafka.Producer {
	kafkaServer := os.Getenv("KAFKA_SERVERS")
	if kafkaServer == "" {
		kafkaServer = "localhost:9092"
	}

	log.Printf("Kafka server: %s", kafkaServer)
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaServer})
	if err != nil {
		panic(err)
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
					// } else {
					// 	fmt.Printf("Delivered to %v, message: %s\n", ev.TopicPartition, string(ev.Value))
				}
			}
		}
	}()

	return p
}
