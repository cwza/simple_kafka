package utils

import (
	"log"
	"os"
	"strconv"
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
