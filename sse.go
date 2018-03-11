package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/gocontrib/log"
	"github.com/gocontrib/pubsub"
)

// GetEventStream streams events in text/event-stream format.
//
// GET /api/event/stream/{channel}
//
func GetEventStream(w http.ResponseWriter, r *http.Request) {
	SendEvents(w, r, getChannels(r))
}

// SendEvents streams events from specified channels as Server Sent Events packets
func SendEvents(w http.ResponseWriter, r *http.Request, channels []string) {
	// make sure that the writer supports flushing
	flusher, ok := w.(http.Flusher)

	if !ok {
		log.Error("current response %T does not implement http.Flusher, plase check your middlewares that wraps response", w)
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	closeNotifier, ok := w.(http.CloseNotifier)
	if !ok {
		log.Error("current response %T does not implement http.CloseNotifier, plase check your middlewares that wraps response", w)
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	sub, err := pubsub.Subscribe(channels)
	if err != nil {
		log.Error("subscribe failed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set the headers related to event streaming
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// channel to send data to client
	var send = make(chan []byte)
	var closeConn = closeNotifier.CloseNotify()
	var heartbeatTicker = time.NewTicker(10 * time.Second) // TODO get from config
	var heatbeat = []byte{}

	var stop = func() {
		log.Info("SSE streaming stopped")
		heartbeatTicker.Stop()
		sub.Close()
		close(send)
	}

	// goroutine to listen all thread events
	go func() {
		if err := recover(); err != nil {
			log.Error("recovered from panic: %+v", err)
			debug.PrintStack()
		}

		log.Info("SSE streaming started")

		for {
			select {
			case msg := <-sub.Read():
				var err error
				var data, ok = msg.([]byte)
				if !ok {
					data, err = json.Marshal(msg)
					if err != nil {
						// TODO should we ignore error messages?
						log.Error("json.Marshal failed with: %+v", err)
						continue
					}
				}
				if len(data) == 0 {
					log.Warning("empty message is not allowed")
					continue
				}
				send <- data
			// listen to connection close and un-register message channel
			case <-sub.CloseNotify():
				log.Info("subscription closed")
				stop()
				return
			case <-closeConn:
				stop()
				return
			case <-heartbeatTicker.C:
				send <- heatbeat
			}
		}
	}()

	for {
		var data, ok = <-send
		if !ok {
			log.Info("connection closed, stop streaming of %v", channels)
			return
		}

		if len(data) == 0 {
			fmt.Fprint(w, ":heartbeat signal\n\n")
		} else {
			fmt.Fprintf(w, "data: %s\n\n", data)
		}

		flusher.Flush()
	}
}

var defaultChannels = []string{"global"}

func getChannels(r *http.Request) []string {
	channel := strings.TrimLeft(chi.URLParam(r, "channel"), "/")
	if len(channel) > 0 {
		return []string{channel}
	}
	s := r.URL.Query().Get("channels")
	if a := split(s); len(a) > 0 {
		return a
	}
	// TODO send events related to current user
	return defaultChannels
}

func split(s string) []string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return []string{}
	}
	var list []string
	for _, c := range strings.Split(s, ",") {
		var t = strings.TrimSpace(c)
		if len(t) > 0 {
			list = append(list, t)
		}
	}
	return list
}
