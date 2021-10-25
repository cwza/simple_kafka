package main

import (
	"fmt"
	"testing"
)

func TestCreateGenValFunc(t *testing.T) {
	startRate := 0
	delta := 6000
	cyclePeriod := 10
	genValFunc := createGenValFunc(startRate, delta, cyclePeriod)

	rates := make([]int, cyclePeriod*10)
	for i := 0; i < cyclePeriod*10; i++ {
		rates[i] = genValFunc()
	}
	fmt.Printf("rates: %v\n", rates)
}
