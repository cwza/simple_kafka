package main

import (
	"fmt"
	"os"
	"testing"
)

func sumInts(as []int) int {
	s := 0
	for _, a := range as {
		s += a
	}
	return s
}

func TestInitConfig(t *testing.T) {
	os.Setenv("CYCLEPERIOD", "5")
	config, err := initConfig("./producer.toml")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("config: %+v\n", config)
}

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
	startRate := 100
	delta := 10
	cyclePeriod := 10
	genMinRateFunc := createGenMinRateFunc(startRate, delta, cyclePeriod)

	genSecRateFunc := createGenSecRateFunc(genMinRateFunc)
	cnt := cyclePeriod * 60 * 2
	rates := make([]int, cnt)
	for i := 0; i < cnt; i++ {
		rates[i] = genSecRateFunc()
	}
	for i := 0; i < cnt/60; i++ {
		tmp := rates[i*60 : (i+1)*60]
		fmt.Printf("%v ", tmp)
		fmt.Printf("%d ", sumInts(tmp))
		fmt.Println()
	}
}
