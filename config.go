package main

import (
	"encoding/gob"
	"fmt"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/spf13/viper"
)

type DBConfig struct {
	Addr string
}

type Config struct {
	Addr string
	DB   *DBConfig
	Nats string
}

var config = &Config{
	Addr: ":3000",
	DB: &DBConfig{
		Addr: "localhost:9080",
	},
	Nats: "nats://nats:4222",
}

func parseConfig() {
	gob.Register(&Event{})
	gob.Register(&api.Assigned{})

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
	config = &Config{
		Addr: viper.GetString("api.addr"),
		DB: &DBConfig{
			Addr: viper.GetString("dgraph.addr"),
		},
		Nats: viper.GetString("nats"),
	}
}
