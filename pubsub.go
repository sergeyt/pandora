package main

import (
	"github.com/gocontrib/pubsub"
)

func startPubsub() {
	// TODO configure nats as main driver in production env
	pubsub.Init()
}

func stopPubsub() {
	pubsub.Cleanup()
}
