package kafkaclient

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// NewConsumer returns a new kafka producer
func NewConsumer(messageHandler func(*kafka.Message)) *kafka.Consumer {
	kafkaServer := os.Getenv("KAFKA_SERVERS")
	if kafkaServer == "" {
		kafkaServer = "localhost:9092"
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        kafkaServer,
		"group.id":                 "rest-api",
		"auto.offset.reset":        "earliest",
		"session.timeout.ms":       6000,
		"go.events.channel.enable": true,
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"request", "response"}, nil)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		run := true

		for run == true {
			select {
			case sig := <-sigchan:
				fmt.Printf("Caught signal %v: terminating\n", sig)
				run = false

			case ev := <-c.Events():
				switch e := ev.(type) {
				case kafka.AssignedPartitions:
					fmt.Fprintf(os.Stderr, "%% %v\n", e)
					c.Assign(e.Partitions)
				case kafka.RevokedPartitions:
					fmt.Fprintf(os.Stderr, "%% %v\n", e)
					c.Unassign()
				case *kafka.Message:
					if messageHandler != nil {
						messageHandler(e)
					}
				case kafka.Error:
					fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
					run = false
				}
			}
		}
	}()

	return c
}
