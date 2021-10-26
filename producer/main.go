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

	startRate   int // msg/min
	delta       int // msg/min
	cyclePeriod int // min
}

var (
	args Args
)

func parseArgs() {
	flag.String(flag.DefaultConfigFlagname, "", "path to config file")
	flag.StringVar(&args.address, "address", "my-cluster-kafka-bootstrap.kafka:9092", "kafka bootstrap address")
	flag.StringVar(&args.topic, "topic", "my-topic", "topic name")
	flag.IntVar(&args.partition, "partition", 8, "number of partitions in this topic")

	flag.IntVar(&args.startRate, "startrate", 0, "msg production rate at start (msg/min)")
	flag.IntVar(&args.delta, "delta", 6000, "increasing decreasing amount each second (msg/min)")
	flag.IntVar(&args.cyclePeriod, "cycleperiod", 10, "duration percycle (min)")

	flag.Parse()
	log.Printf("args: %+v\n", args)
}

func createGenMinRateFunc(start int, delta int, cyclePeriod int) func() int {
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

func createGenSecRateFunc(genMinRateFunc func() int) func() int {
	var minRate int
	sec := 0
	return func() int {
		if sec >= 60 {
			sec = 0
		}
		if sec == 0 {
			minRate = genMinRateFunc()
		}
		secRate := minRate / 60
		sec++
		return secRate
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

func run(writer *kafka.Writer) {
	genMinRateFunc := createGenMinRateFunc(args.startRate, args.delta, args.cyclePeriod)
	genSecRateFunc := createGenSecRateFunc(genMinRateFunc)
	for {
		cnt := genSecRateFunc()
		err := send(writer, cnt)
		if err != nil {
			log.Printf("WARNING: failed to send %d messages, %s\n", cnt, err)
		}
		log.Printf("send %d msgs\n", cnt)
		time.Sleep(time.Second)
	}
}

func main() {
	parseArgs()

	err := utils.CreateTopic(args.address, args.topic, args.partition)
	if err != nil {
		log.Fatalf("failed to create topic, %s", err)
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(args.address),
		Topic:    args.topic,
		Balancer: &kafka.LeastBytes{},
	}
	run(writer)
}
