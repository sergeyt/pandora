package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Addr string
}

var config Config = Config{}

func parseConfig() {
	viper.SetConfigType("toml")
	viper.SetConfigName("pandora")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("bad config file: %s", err))
	}

	initConfig()
}

func initConfig() {
	config = Config{
		Addr: viper.GetString("http.addr"),
	}
}
