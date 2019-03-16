package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Addr string
}

type ElasticSearchConfig struct {
	URL string
}

var (
	ServerAddr = ":3000"
	DB         = &DBConfig{
		Addr: "localhost:9080",
	}
	ElasticSearch = &ElasticSearchConfig{
		URL: "http://elasticsearch:9200",
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
		URL: viper.GetString("elasticsearch.url"),
	}
	Nats = viper.GetString("nats")
}

func ServerURL() string {
	s := os.Getenv("SERVER_URL")
	if len(s) > 0 {
		return s
	}
	hostname, err := os.Hostname()
	if err == nil {
		return fmt.Sprintf("http://%s", hostname)
	}
	return "http://localhost:4200"
}
