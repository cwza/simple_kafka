package main

import (
	"fmt"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func createTopic(address string, topic string, partitionCnt int) error {
	conn, err := kafka.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to dial, %s", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller, %s", err)
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return fmt.Errorf("failed to dial with controller, %s", err)
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafka.TopicConfig{Topic: topic, NumPartitions: partitionCnt, ReplicationFactor: 1})
	if err != nil {
		return fmt.Errorf("failed to create topics, %s", err)
	}
	return nil
}
