package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gocontrib/pubsub"
	_ "github.com/gocontrib/pubsub/nats"
	"github.com/spf13/viper"
)

func main() {
	parseConfig()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		stop()
		initConfig()
		go start()
	})

	die := make(chan bool)
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	go func() {
		<-sig
		die <- true
	}()

	go start()
	<-die

	stop()
}

func start() {
	startHub()
	startServer()
}

func stop() {
	pubsub.Cleanup()
	stopServer()
}

func startHub() {
	nats := config.Nats
	for attemt := 0; attemt < 30; attemt = attemt + 1 {
		err := pubsub.Init(pubsub.HubConfig{
			"driver": "nats",
			"url":    nats,
		})
		if err == nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
	log.Fatalf("cannot initialize hub")
}
