package main

import (
	"github.com/gocontrib/pubsub"
)

func startPubsub() {
	// TODO configure nats as main driver at production
	pubsub.Init()
}

func stopPubsub() {
	pubsub.Cleanup()
}
