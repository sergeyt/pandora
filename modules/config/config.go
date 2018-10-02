package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Addr string
}

type ElasticSearchConfig struct {
	URL       string
	IndexName string
}

var (
	ServerAddr = ":3000"
	DB         = &DBConfig{
		Addr: "localhost:9080",
	}
	ElasticSearch = &ElasticSearchConfig{
		URL:       "http://elasticsearch:9200",
		IndexName: "pandora_data",
	}
	Nats = "nats://nats:4222"
)

func Parse() {
	viper.SetConfigType("toml")
	viper.SetConfigName("pandora")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("bad config file: %s", err))
	}

	Init()
}

func Init() {
	ServerAddr = viper.GetString("api.addr")
	DB = &DBConfig{
		Addr: viper.GetString("dgraph.addr"),
	}
	ElasticSearch = &ElasticSearchConfig{
		URL:       viper.GetString("elasticsearch.url"),
		IndexName: viper.GetString("elasticsearch.index"),
	}
	Nats = viper.GetString("nats")
}
