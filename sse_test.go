package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gocontrib/rest"
	"github.com/stretchr/testify/assert"
)

func TestEventStream(t *testing.T) {
	port := os.Getenv("HTTP_PORT")
	host := fmt.Sprintf("http://localhost:%s", port)
	client := rest.NewClient(rest.Config{
		BaseURL:       host,
		Authorization: "local_admin",
		Verbose:       true,
	})

	events := make(chan *rest.Event, 1)
	die := make(chan bool)
	timer := time.NewTimer(5 * time.Second)

	go func() {
		fmt.Println("READING EVENT STREAM")
		err := client.EventStream("/api/event/stream", events)
		assert.Nil(t, err)
		<-die
		close(events)
	}()

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("POST DATA")
		var result map[string]interface{}
		err := client.Post("/api/data/user", &TestUser{
			Name: "bob",
			Age:  21,
		}, &result)
		assert.Nil(t, err)
	}()

	select {
	case e := <-events:
		assert.NotNil(t, e)
	case <-timer.C:
		assert.Fail(t, "timeout")
	}

	die <- true
}
