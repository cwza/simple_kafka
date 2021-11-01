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
	var secRates []int
	sec := 0
	return func() int {
		if sec >= 60 {
			sec = 0
		}
		if sec == 0 {
			secRates = make([]int, 60)
			minRate := genMinRateFunc()
			secRate := minRate / 60
			remain := minRate % 60
			for i := 0; i < 60; i++ {
				secRates[i] = secRate
			}
			for i := 0; i < remain; i++ {
				secRates[i]++
			}
		}
		secRate := secRates[sec]
		sec++
		return secRate
	}
}
