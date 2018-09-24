package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/go-chi/chi"
	"github.com/gocontrib/pubsub"
)

func dataAPI(r chi.Router) {
	r = r.With(transactionMiddleware)

	r.Post("/api/data/query", queryHandler)
	r.Get("/api/data/{type}/{id}", readHandler)

	// mutation api
	r.Post("/api/data/{type}", mutateHandler)
	r.Put("/api/data/{type}/{id}", mutateHandler)
	r.Delete("/api/data/{type}/{id}", deleteHandler)

	// TODO allow to delete triples from graph
	// TODO consider to expose raw api for admin users

	// edges api
}

func transactionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := newDgraphClient()
		if err != nil {
			sendError(w, err)
			return
		}

		tx := c.NewTxn()
		defer tx.Discard(r.Context())

		ctx := context.WithValue(r.Context(), "tx", tx)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func transaction(r *http.Request) *dgo.Txn {
	return r.Context().Value("tx").(*dgo.Txn)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	query, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendError(w, err)
		return
	}

	tx := transaction(r)

	resp, err := tx.Query(r.Context(), string(query))
	if err != nil {
		sendError(w, err)
		return
	}

	w.Header().Set("Content-Type", TypeJSON)
	w.Write(resp.GetJson())
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	tx := transaction(r)
	id := chi.URLParam(r, "id")

	query := fmt.Sprintf(`{
  q(func: uid(%s)) {
    uid
    expand(_all_) {
      expand(_all_)
	}
  }
}`, id)
	resp, err := tx.Query(r.Context(), query)
	if err != nil {
		sendError(w, err)
		return
	}

	var result struct {
		Results []map[string]interface{} `json:"q"`
	}
	err = json.Unmarshal(resp.GetJson(), &result)
	if err != nil {
		sendError(w, err)
		return
	}
	if len(result.Results) == 0 {
		sendError(w, fmt.Errorf("not found"))
		return
	}
	sendJSON(w, result.Results[0])
}

func mutateHandler(w http.ResponseWriter, r *http.Request) {
	tx := transaction(r)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendError(w, err)
		return
	}

	resp, err := tx.Mutate(r.Context(), &api.Mutation{
		SetJson:   data,
		CommitNow: true,
	})
	if err != nil {
		sendError(w, err)
		return
	}

	err = sendJSON(w, resp)
	if err != nil {
		return
	}

	resourceType := chi.URLParam(r, "type")
	id := chi.URLParam(r, "id")

	// TODO set CreatedBy
	sendEvent(&Event{
		Action:       r.Method,
		ResourceID:   id,
		ResourceType: resourceType,
		CreatedAt:    time.Now(),
		DbResponse:   resp,
	})
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	tx := transaction(r)
	resourceType := chi.URLParam(r, "type")
	id := chi.URLParam(r, "id")

	resp, err := tx.Mutate(r.Context(), &api.Mutation{
		DelNquads: []byte("<" + id + "> * * .\n"),
		CommitNow: true,
	})
	if err != nil {
		sendError(w, err)
		return
	}

	err = sendJSON(w, resp)
	if err != nil {
		return
	}

	// TODO set CreatedBy
	sendEvent(&Event{
		Action:       r.Method,
		ResourceID:   id,
		ResourceType: resourceType,
		CreatedAt:    time.Now(),
		DbResponse:   resp,
	})
}

// TODO also implement persistence of events for some period of time
func sendEvent(evt *Event) {
	go func() {
		chans := []string{
			"global",
			fmt.Sprintf("%s/%s", evt.ResourceType, evt.ResourceID),
			// TODO push to user channel too
		}
		pubsub.Publish(chans, evt)
	}()
}

type Event struct {
	Action       string      `json:"action"`
	ResourceID   string      `json:"resource_id"`   // resource id
	ResourceType string      `json:"resource_type"` // resource type
	CreatedBy    string      `json:"created_by"`
	CreatedAt    time.Time   `json:"created_at"`
	DbResponse   interface{} `json:"db_response"`
}
