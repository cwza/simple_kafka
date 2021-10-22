package main

import (
	"context"
	"fmt"
	"log"

	"github.com/cwza/simple_kafka/utils"
	"github.com/segmentio/kafka-go"
)

var (
	reader *kafka.Reader
)

func init() {
	address := utils.GetEnvStr("ADDRESS")
	topic := utils.GetEnvStr("TOPIC")
	partitionCnt := utils.GetEnvInt("PARTITION_CNT")
	groupId := utils.GetEnvStr("GROUP_ID")

	err := utils.CreateTopic(address, topic, partitionCnt)
	if err != nil {
		log.Fatalf("failed to create topic: %s", err)
	}

	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{address},
		GroupID:  groupId,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func main() {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("failed to read msg: %s", err)
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	// if err := reader.Close(); err != nil {
	// 	log.Fatal("failed to close reader:", err)
	// }
}
