package consumer

import (
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func Reader() *kafka.Reader {
	kafkaBroker := os.Getenv("KAFKA_BROKER")

	if kafkaBroker == "" {
		panic("Unable to find the Kafka broker address")
	}

	readerConfig := kafka.ReaderConfig{
		Brokers:     []string{kafkaBroker},
		Topic:       "tips.created",
		Partition:   0,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		MaxWait:     1 * time.Second,
		GroupID:     "mail-consumer",
		StartOffset: kafka.LastOffset,
	}
	reader := kafka.NewReader(readerConfig)

	return reader
}
