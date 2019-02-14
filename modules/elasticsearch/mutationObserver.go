package elasticsearch

import (
	"encoding/json"
	"time"

	"github.com/gocontrib/pubsub"
	log "github.com/sirupsen/logrus"
)

// MutationObserver runs monitor of all mutation events to replicate data
func MutationObserver(restart chan bool) {
	ch, err := pubsub.Subscribe([]string{"global"})
	if err != nil {
		log.Info("elasticseach: cannot subscribe on global channel: %v\n", err)
		log.Info("elasticseach: retry after one second")
		time.Sleep(1 * time.Second)
		go MutationObserver(restart)
		return
	}

	log.Info("elasticseach: mutation observer started")

	for {
		select {
		case msg := <-ch.Read():
			go mutate(msg)
		case <-restart:
			go MutationObserver(restart)
			return
		case <-ch.CloseNotify():
			time.Sleep(1 * time.Second)
			log.Info("this subscription is closed")
			log.Info("retry after one second")
			go MutationObserver(restart)
			return
		}
	}
}

func mutate(msg interface{}) {
	b, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("elasticsearch: json encoding error: %v\n", err)
		return
	}
	var event pubsub.Event
	err = json.Unmarshal(b, &event)
	if err != nil {
		log.Errorf("elasticsearch: json decoding error: %v\n", err)
		return
	}

	if event.Result != nil {
		c := makeClient()
		idx := "pandora_" + event.ResourceType
		c.Push(idx, event.Result, event.ResourceID)
	}
}
