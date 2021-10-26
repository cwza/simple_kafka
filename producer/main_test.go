package main

import (
	"fmt"
	"testing"
)

func TestCreateGenMinRateFunc(t *testing.T) {
	startRate := 0
	delta := 6000
	cyclePeriod := 10
	genMinRateFunc := createGenMinRateFunc(startRate, delta, cyclePeriod)

	cnt := cyclePeriod * 10
	rates := make([]int, cnt)
	for i := 0; i < cnt; i++ {
		rates[i] = genMinRateFunc()
	}
	fmt.Printf("rates: %v\n", rates)
}

func TestCreateGenSecRateFunc(t *testing.T) {
	startRate := 0
	delta := 6000
	cyclePeriod := 10
	genMinRateFunc := createGenMinRateFunc(startRate, delta, cyclePeriod)

	genSecRateFunc := createGenSecRateFunc(genMinRateFunc)
	cnt := cyclePeriod * 60 * 2
	rates := make([]int, cnt)
	for i := 0; i < cnt; i++ {
		rates[i] = genSecRateFunc()
	}
	fmt.Printf("rates: %v\n", rates)
}
