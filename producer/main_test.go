package main

import (
	"fmt"
	"testing"
)

func TestCreateGenValFunc(t *testing.T) {
	startRate := 6000
	delta := 8000
	cyclePeriod := 10
	genValFunc := createGenValFunc(startRate, delta, cyclePeriod)

	rates := make([]int, cyclePeriod)
	for i := 0; i < cyclePeriod; i++ {
		rates[i] = genValFunc()
	}
	fmt.Printf("rates: %v\n", rates)
}
