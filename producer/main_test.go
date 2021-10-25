package main

import (
	"fmt"
	"testing"
)

func sumInts(values []int) int {
	s := 0
	for _, value := range values {
		s += value
	}
	return s
}

func TestCreateGenValFunc(t *testing.T) {
	startRate := 10
	delta := 8
	cyclePeriod := 600
	genValFunc := createGenValFunc(startRate, delta, cyclePeriod)

	secRates := make([]int, cyclePeriod)
	for i := 0; i < cyclePeriod; i++ {
		secRates[i] = genValFunc()
	}
	fmt.Printf("secRates: %v\n", secRates)

	minRates := make([]int, 0)
	for i := 0; i < cyclePeriod; i += 60 {
		start := i
		end := i + 60
		if end > cyclePeriod {
			end = cyclePeriod
		}
		minRate := sumInts(secRates[start:end])
		minRates = append(minRates, minRate)
	}
	fmt.Printf("minRates: %v\n", minRates)
}
