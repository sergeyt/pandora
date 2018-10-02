package elasticsearch

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gocontrib/pubsub"
	"github.com/sergeyt/pandora/modules/config"
)

// MutationObserver runs monitor of all mutation events to replicate data
func MutationObserver(restart chan bool) {
	ch, err := pubsub.Subscribe([]string{"global"})
	if err != nil {
		log.Printf("cannot subscribe on global channel: %v", err)
		log.Println("retry after one second")
		time.Sleep(1 * time.Second)
		return
	}

	for {
		select {
		case msg := <-ch.Read():
			go mutate(msg)
		case <-restart:
			return
		case <-ch.CloseNotify():
			time.Sleep(1 * time.Second)
			log.Println("this subscription is closed")
			log.Println("retry after one second")
			go MutationObserver(restart)
			return
		}
	}
}

func mutate(msg interface{}) {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Printf("elasticsearch process: json encoding error: %v", err)
		return
	}
	var event pubsub.Event
	err = json.Unmarshal(b, &event)
	if err != nil {
		log.Printf("elasticsearch process: json decoding error: %v", err)
		return
	}

	if event.Result != nil {
		c := makeClient()
		c.Push(config.ElasticSearch.IndexName, event.Result)
	}
}
