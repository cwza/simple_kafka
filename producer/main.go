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

	startRate   int // msg/sec
	delta       int // msg/sec
	cyclePeriod int // sec
}

var (
	args Args
)

func parseArgs() {
	flag.String(flag.DefaultConfigFlagname, "", "path to config file")
	flag.StringVar(&args.address, "address", "my-cluster-kafka-bootstrap.kafka:9092", "kafka bootstrap address")
	flag.StringVar(&args.topic, "topic", "my-topic", "topic name")
	flag.IntVar(&args.partition, "partition", 1, "number of partitions in this topic")

	flag.IntVar(&args.startRate, "startrate", 1000, "msg production rate at start (msg/sec)")
	flag.IntVar(&args.delta, "delta", 16, "increasing decreasing amount each second (msg/sec)")
	flag.IntVar(&args.cyclePeriod, "cycleperiod", 600, "duration percycle (sec)")

	flag.Parse()
	log.Printf("args: %+v\n", args)
}

func createGenValFunc(start int, delta int, cyclePeriod int) func() int {
	val := start
	i := 0
	return func() int {
		if i < cyclePeriod/2 {
			val = val + delta
		} else {
			val = val - delta
		}
		i++
		if i >= cyclePeriod {
			i = 0
		}
		return val
	}
}

func send(writer *kafka.Writer, cnt int) error {
	if cnt < 0 {
		return fmt.Errorf("cnt is zero")
	}
	msgs := make([]kafka.Message, cnt)
	for i := 0; i < cnt; i++ {
		msgs[i] = kafka.Message{Key: []byte("cnt"), Value: []byte(fmt.Sprintf("%d", cnt))}
	}
	err := writer.WriteMessages(context.Background(), msgs...)
	return err
}

func run() {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(args.address),
		Topic:    args.topic,
		Balancer: &kafka.LeastBytes{},
	}

	genValFunc := createGenValFunc(args.startRate, args.delta, args.cyclePeriod)
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		cnt := genValFunc()
		err := send(writer, cnt)
		if err != nil {
			fmt.Printf("WARNING: failed to write %d messages, %s\n", cnt, err)
		}
		fmt.Printf("send rate %d msg/sec\n", cnt)
	}
}

func main() {
	parseArgs()

	err := utils.CreateTopic(args.address, args.topic, args.partition)
	if err != nil {
		log.Fatalf("failed to create topic, %s", err)
	}

	run()
}
