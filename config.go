package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Addr string
}

type Config struct {
	Addr string
	DB   DBConfig
}

var config Config = Config{
	Addr: ":8888",
	DB: DBConfig{
		Addr: "localhost:9080",
	},
}

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
		DB: DBConfig{
			Addr: viper.GetString("dgraph.addr"),
		},
	}
}
