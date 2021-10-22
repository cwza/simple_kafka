package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cwza/simple_kafka/utils"
	"github.com/namsral/flag"
	"github.com/segmentio/kafka-go"
)

type Args struct {
	address   string
	topic     string
	partition int
	interval  int
}

var (
	args   Args
	writer *kafka.Writer
)

func parseArgs() {
	flag.String(flag.DefaultConfigFlagname, "", "path to config file")
	flag.StringVar(&args.address, "address", "my-cluster-kafka-bootstrap.kafka:9092", "kafka bootstrap address")
	flag.StringVar(&args.topic, "topic", "my-topic", "topic name")
	flag.IntVar(&args.partition, "partition", 1, "number of partitions in this topic")
	flag.IntVar(&args.interval, "interval", 30000, "")
	flag.Parse()
	log.Printf("args: %+v\n", args)
}

func init() {
	parseArgs()

	err := utils.CreateTopic(args.address, args.topic, args.partition)
	if err != nil {
		log.Fatalf("failed to create topic, %s", err)
	}

	writer = &kafka.Writer{
		Addr:     kafka.TCP(args.address),
		Topic:    args.topic,
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
	ticker := time.NewTicker(time.Duration(args.interval) * time.Millisecond)
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
