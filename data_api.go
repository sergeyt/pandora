package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/go-chi/chi"
	"github.com/gocontrib/auth"
	"github.com/gocontrib/pubsub"
)

func dataAPI(r chi.Router) {
	r = r.With(authMiddleware)
	r = r.With(transactionMiddleware)

	r.Post("/api/query", queryHandler)
	r.Get("/api/me", meHandler)
	r.Get("/api/data/{type}/list", listHandler)
	r.Get("/api/data/{type}/{id}", readHandler)

	// mutation api
	r.Post("/api/data/{type}", mutateHandler)
	r.Put("/api/data/{type}/{id}", mutateHandler)
	r.Delete("/api/data/{type}/{id}", deleteHandler)

	// TODO allow to delete triples from graph
	// TODO consider to expose raw api for admin users

	// edges api
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

func listHandler(w http.ResponseWriter, r *http.Request) {
	resourceType := strings.ToLower(chi.URLParam(r, "type"))
	pg, err := parsePagination(r)
	if err != nil {
		sendError(w, err)
		return
	}

	tx := transaction(r)

	data, err := readList(r.Context(), tx, nodeLabel(resourceType), pg)
	if err != nil {
		sendError(w, err)
		return
	}

	sendJSON(w, data)
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	user := auth.GetRequestUser(r)
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	readHandlerByID(w, r, user.GetID())
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	readHandlerByID(w, r, id)
}

func readHandlerByID(w http.ResponseWriter, r *http.Request, id string) {
	tx := transaction(r)

	data, err := readNode(r.Context(), tx, id)
	if err != nil {
		sendError(w, err)
		return
	}

	sendJSON(w, data)
}

func mutateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceType := strings.ToLower(chi.URLParam(r, "type"))
	id := chi.URLParam(r, "id")
	nodeLabel := nodeLabel(resourceType)
	user := auth.GetContextUser(ctx)

	var in OrderedJSON
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		sendError(w, err)
		return
	}

	isNew := len(id) == 0
	tx := transaction(r)
	now := time.Now()

	in["modified_at"] = now
	in["modified_by"] = user.GetID()

	if isNew {
		in[nodeLabel] = ""
		in["created_at"] = now
		in["created_by"] = user.GetID()
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

	if isNew {
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
			if len(results) == 1 {
				id = uid
			}
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

	sendEvent(user, &pubsub.Event{
		Action:       r.Method,
		Method:       r.Method,
		URL:          r.URL.String(),
		ResourceID:   id,
		ResourceType: resourceType,
		Payload:      &in,
		CreatedBy:    user.GetID(),
		CreatedAt:    time.Now(),
		Result:       out,
	})
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	resourceType := strings.ToLower(chi.URLParam(r, "type"))
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	user := auth.GetContextUser(ctx)

	tx := transaction(r)

	resp, err := tx.Mutate(ctx, &api.Mutation{
		DelNquads: []byte("<" + id + "> * * .\n"),
		CommitNow: true,
	})
	if err != nil {
		sendError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)

	sendEvent(user, &pubsub.Event{
		Action:       r.Method,
		Method:       r.Method,
		URL:          r.URL.String(),
		ResourceID:   id,
		ResourceType: resourceType,
		CreatedBy:    user.GetID(),
		CreatedAt:    time.Now(),
		Result:       resp,
	})
}

// TODO later implement persistence of events for some period of time
func sendEvent(user auth.User, evt *pubsub.Event) {
	go func() {
		chans := []string{
			"global",
			fmt.Sprintf("%s/%s", evt.ResourceType, evt.ResourceID),
			fmt.Sprintf("user/%s", user.GetID()),
		}
		pubsub.Publish(chans, evt)
	}()
}
