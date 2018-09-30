package main

import (
	"testing"
	"time"

	"github.com/gocontrib/rest"
	"github.com/r3labs/sse"
	"github.com/stretchr/testify/assert"
)

func TestEventStream(t *testing.T) {
	host := "http://localhost:4200"
	client := rest.NewClient(rest.Config{
		BaseURL:       host,
		Authorization: "local_admin",
		Verbose:       true,
	})

	stream := sse.NewClient(host + "/api/event/stream")
	events := make(chan *sse.Event)

	go func() {
		err := stream.SubscribeChan("", events)
		assert.Nil(t, err)
	}()
	defer stream.Unsubscribe(events)

	go func() {
		var result map[string]interface{}
		err := client.Post("/api/data/user", &TestUser{
			Name: "bob",
			Age:  21,
		}, &result)
		assert.Nil(t, err)
	}()

	for {
		select {
		case e := <-events:
			assert.NotNil(t, e)
			return
		case <-time.After(1 * time.Second):
			assert.Fail(t, "timeout")
		}
	}
}
