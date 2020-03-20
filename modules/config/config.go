package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sergeyt/pandora/modules/env"
)

type DgraphConfig struct {
	RpcURL  string // e.g. dgraph:9080
	HttpURL string // e.g. http://dgraph:8080
}

type ElasticSearchConfig struct {
	URL string
}

var (
	ServerPort = env.Get("PORT", "3000")
	DGraph     = &DgraphConfig{
		RpcURL:  env.Get("DGRAPH_RPC_URL", "localhost:9080"),
		HttpURL: env.Get("DGRAPH_HTTP_URL", "http://localhost:8080"),
	}
	ElasticSearch = &ElasticSearchConfig{
		URL: env.Get("ES_HOSTS", "http://localhost:9200"),
	}
	NatsURL = env.Get("NATS_URI", "nats://localhost:4222")
)

func Reload() {
	ServerPort = env.Get("PORT", "3000")
	DGraph = &DgraphConfig{
		RpcURL:  env.Get("DGRAPH_RPC_URL", "localhost:9080"),
		HttpURL: env.Get("DGRAPH_HTTP_URL", "http://localhost:8080"),
	}
	ElasticSearch = &ElasticSearchConfig{
		URL: env.Get("ES_HOSTS", "http://localhost:9200"),
	}
	NatsURL = env.Get("NATS_URI", "nats://localhost:4222")
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
	// TODO detect https is enabled
	secure := true
	scheme := "http"
	portVar := "HTTP_PORT"
	if secure {
		scheme = "https"
		portVar = "HTTPS_PORT"
	}
	if port, err := strconv.ParseInt(portVar, 10, 64); err == nil {
		if port == 80 {
			return fmt.Sprintf("%s://%s", scheme, hostname)
		}
		return fmt.Sprintf("%s://%s:%d", scheme, hostname, port)
	}
	return fmt.Sprintf("%s://%s", scheme, hostname)
}
