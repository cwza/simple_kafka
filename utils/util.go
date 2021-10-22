package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func GetEnvStr(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Env of %s is empty\n", key)
	}
	return val
}

func GetEnvInt(key string) int {
	val := GetEnvStr(key)
	ans, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("Env of %s is not a int\n", key)
	}
	return ans
}

func CreateTopic(address string, topic string, partitionCnt int) error {
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
