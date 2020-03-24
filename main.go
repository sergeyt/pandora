package main

import (
	"os"
	"os/signal"

	"github.com/gocontrib/pubsub"
	_ "github.com/gocontrib/pubsub/nats"
	"github.com/sergeyt/pandora/modules/auth"
	"github.com/sergeyt/pandora/modules/cloudstore"
	"github.com/sergeyt/pandora/modules/config"
	"github.com/sergeyt/pandora/modules/dgraph"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetOutput(os.Stdout)

	restart := make(chan bool)

	die := make(chan bool)
	sig := make(chan os.Signal)

	signal.Notify(sig, os.Interrupt, os.Kill)
	go func() {
		<-sig
		die <- true
	}()

	go start(restart)
	<-die

	stop()
}

func start(restart chan bool) {
	cloudstore.InitStore()
	dgraph.InitSchema()
	auth.InitUsers()

	startHub()
	// go elasticsearch.MutationObserver(restart)
	startServer()
}

func stop() {
	pubsub.Cleanup()
	stopServer()
}

func startHub() {
	conf := pubsub.HubConfig{
		"driver": "nats",
		"url":    config.NatsURL,
	}
	err := pubsub.Init(conf)
	if err != nil {
		log.Fatalf("cannot initialize hub")
	}
}
