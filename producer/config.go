package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Address   string `mapstructure:"address"`
	Topic     string `mapstructure:"topic"`
	Partition int    `mapstructure:"partition"`

	StartRate   int `mapstructure:"startrate"`
	Delta       int `mapstructure:"delta"`
	CyclePeriod int `mapstructure:"cycleperiod"`
}

func initConfig(filepath string) (Config, error) {
	viper.SetConfigFile(filepath)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	config := Config{}
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("read config failed: %s", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("unmarshal config failed: %s", err)
	}

	if config.Address == "" {
		return Config{}, fmt.Errorf("config.Address is invalid")
	}
	if config.Topic == "" {
		return Config{}, fmt.Errorf("config.Topic is invalid")
	}
	if config.Partition <= 0 {
		return Config{}, fmt.Errorf("config.Partition is invalid")
	}
	if config.StartRate < 0 {
		return Config{}, fmt.Errorf("config.StartRate is invalid")
	}
	if config.CyclePeriod <= 0 {
		return Config{}, fmt.Errorf("config.CyclePeriod is invalid")
	}

	return config, nil
}
