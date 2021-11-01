package main

import (
	"context"
	"flag"
	"log"
	"runtime"

	"github.com/segmentio/kafka-go"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "cfgpath", "./consumer-config.toml", "config file path")
}

func run(reader *kafka.Reader) {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("WARNING: failed to read msg, %s\n", err)
		}
		log.Printf("read msg from partition: %d\n", m.Partition)
	}
}

func main() {
	flag.Parse()

	config, err := initConfig(configPath)
	if err != nil {
		log.Fatalf("failed to init config, %s", err)
	}

	err = createTopic(config.Address, config.Topic, config.Partition)
	if err != nil {
		log.Fatalf("failed to create topic, %s", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.Address},
		GroupID:  config.GroupId,
		Topic:    config.Topic,
		MinBytes: 10e6, // 10MB
		MaxBytes: 50e6, // 50MB
	})
	// run(reader)
	for i := 0; i < runtime.NumCPU(); i++ {
		go run(reader)
	}
	select {}
}
