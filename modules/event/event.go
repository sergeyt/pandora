package event

import (
	"fmt"

	"github.com/gocontrib/auth"
	"github.com/gocontrib/pubsub"
	log "github.com/sirupsen/logrus"
)

// TODO later implement persistence of events for some period of time
func Send(user auth.User, evt *pubsub.Event) {
	if evt == nil {
		panic("event is required")
	}
	go func() {
		chans := []string{
			"global",
			fmt.Sprintf("%s/%s", evt.ResourceType, evt.ResourceID),
		}
		if user != nil && user.GetID() != "" {
			uchan := fmt.Sprintf("user/%s", user.GetID())
			chans = append(chans, uchan)
		}
		err := pubsub.Publish(chans, evt)
		if err != nil {
			log.Errorf("pubsub.Publish fail: %v", err)
		}

		// TODO notify other users who subscribed to changed resource
	}()
}
