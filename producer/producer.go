package main

import (
	"context"
	"log"
	"time"

	"github.com/cwza/simple_kafka/utils"
	"github.com/segmentio/kafka-go"
)

var (
	interval int
	writer   *kafka.Writer
)

func init() {
	address := utils.GetEnvStr("ADDRESS")
	topic := utils.GetEnvStr("TOPIC")
	interval = utils.GetEnvInt("INTERVAL")
	writer = &kafka.Writer{
		Addr:     kafka.TCP(address),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func send(key string, val string) error {
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: []byte(val),
	})
	return err
}

func main() {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	for range ticker.C {
		err := send("keykey", "valval")
		if err != nil {
			log.Fatalf("failed to write messages: %s", err)
		}
	}
	// if err := writer.Close(); err != nil {
	// 	log.Fatal("failed to close writer:", err)
	// }
}
