package main

import (
	"context"
	"fmt"
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
	partitionCnt := utils.GetEnvInt("PARTITION_CNT")
	interval = utils.GetEnvInt("INTERVAL")

	err := utils.CreateTopic(address, topic, partitionCnt)
	if err != nil {
		log.Fatalf("failed to create topic, %s", err)
	}

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
			fmt.Printf("WARNING: failed to write messages, %s", err)
		}
		fmt.Println("send keykey:valval")
	}
	// if err := writer.Close(); err != nil {
	// 	log.Fatal("failed to close writer:", err)
	// }
}
