package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "cfgpath", "./producer.toml", "config file path")
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

func run(writer *kafka.Writer, genSecRateFunc func() int) {
	for range time.Tick(time.Second) {
		cnt := genSecRateFunc()
		err := send(writer, cnt)
		if err != nil {
			log.Printf("WARNING: failed to send %d messages, %s\n", cnt, err)
		}
		log.Printf("send %d msgs\n", cnt)
	}
}

func main() {
	flag.Parse()

	config, err := initConfig(configPath)
	if err != nil {
		log.Fatalf("failed to init config, %s", err)
	}
	log.Printf("config: %+v\n", config)

	err = createTopic(config.Address, config.Topic, config.Partition)
	if err != nil {
		log.Fatalf("failed to create topic, %s", err)
	}

	genSecRateFunc := createGenSecRateFunc(createGenMinRateFunc(config.Rates, config.Cnts))
	writer := &kafka.Writer{
		Addr:     kafka.TCP(config.Address),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	}
	run(writer, genSecRateFunc)
}
