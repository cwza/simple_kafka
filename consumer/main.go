package main

import (
	"context"
	"log"

	"github.com/cwza/simple_kafka/utils"
	"github.com/namsral/flag"
	"github.com/segmentio/kafka-go"
)

type Args struct {
	address   string
	topic     string
	partition int
	groupid   string
}

var (
	args   Args
	reader *kafka.Reader
)

func parseArgs() {
	flag.String(flag.DefaultConfigFlagname, "", "path to config file")
	flag.StringVar(&args.address, "address", "my-cluster-kafka-bootstrap.kafka:9092", "kafka bootstrap address")
	flag.StringVar(&args.topic, "topic", "my-topic", "topic name")
	flag.IntVar(&args.partition, "partition", 1, "number of partitions in this topic")
	flag.StringVar(&args.groupid, "groupid", "simple-kafka-consumer", "consumer group id")
	flag.Parse()
	log.Printf("args: %+v\n", args)
}

func init() {
	parseArgs()

	err := utils.CreateTopic(args.address, args.topic, args.partition)
	if err != nil {
		log.Fatalf("failed to create topic, %s", err)
	}

	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{args.address},
		GroupID:  args.groupid,
		Topic:    args.topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func main() {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("WARNING: failed to read msg, %s", err)
		}
		log.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	// if err := reader.Close(); err != nil {
	// 	log.Fatal("failed to close reader:", err)
	// }
}
