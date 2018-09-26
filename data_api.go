package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/go-chi/chi"
	"github.com/gocontrib/pubsub"
)

func dataAPI(r chi.Router) {
	r = r.With(authMiddleware)
	r = r.With(transactionMiddleware)

	r.Post("/api/query", queryHandler)
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

	data, err := readNode(r.Context(), tx, id)
	if err != nil {
		sendError(w, err)
		return
	}

	sendJSON(w, data)
}

func mutateHandler(w http.ResponseWriter, r *http.Request) {
	resourceType := strings.ToLower(chi.URLParam(r, "type"))
	id := chi.URLParam(r, "id")
	nodeLabel := "_" + resourceType

	var in OrderedJSON
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		sendError(w, err)
		return
	}

	ctx := r.Context()
	tx := transaction(r)
	now := time.Now()

	in["modified_at"] = now

	if len(id) == 0 {
		in[nodeLabel] = ""
		in["created_at"] = now
	} else {
		in["uid"] = id
	}

	data, err := in.ToJSON("uid", nodeLabel)
	if err != nil {
		sendError(w, err)
		return
	}

	resp, err := tx.Mutate(ctx, &api.Mutation{
		SetJson: data,
	})
	if err != nil {
		sendError(w, err)
		return
	}

	var results []map[string]interface{}

	if len(id) == 0 {
		results = make([]map[string]interface{}, len(resp.Uids))
		i := 0
		for _, uid := range resp.Uids {
			result, err := readNode(ctx, tx, uid)
			if err != nil {
				sendError(w, err)
				return
			}
			results[i] = result
			i = i + 1
		}
	} else {
		result, err := readNode(ctx, tx, id)
		if err != nil {
			sendError(w, err)
			return
		}
		results = []map[string]interface{}{result}
	}

	err = tx.Commit(ctx)
	if err != nil {
		sendError(w, err)
		return
	}

	var out interface{} = results
	if len(results) == 1 {
		out = results[0]
	}

	err = sendJSON(w, out)
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

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	resourceType := strings.ToLower(chi.URLParam(r, "type"))
	id := chi.URLParam(r, "id")

	tx := transaction(r)

	resp, err := tx.Mutate(r.Context(), &api.Mutation{
		DelNquads: []byte("<" + id + "> * * .\n"),
		CommitNow: true,
	})
	if err != nil {
		sendError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)

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
