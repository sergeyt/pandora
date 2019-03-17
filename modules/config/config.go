package config

import (
	"fmt"
	"os"
	"strconv"

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

// Hostname reads HOSTNAME env var or os.Hostname used for your app
func Hostname() string {
	hostname := os.Getenv("HOSTNAME")
	if len(hostname) > 0 {
		return hostname
	}
	hostname, err := os.Hostname()
	if err == nil {
		return hostname
	}
	return "localhost"
}

// ServerURL returns base URL including hostname and port
func ServerURL() string {
	hostname := Hostname()
	if port, err := strconv.ParseInt(os.Getenv("HTTP_PORT"), 10, 64); err == nil {
		if port == 80 {
			return fmt.Sprintf("http://%s", hostname)
		}
		return fmt.Sprintf("http://%s:%d", hostname, port)
	}
	return fmt.Sprintf("http://%s", hostname)
}
