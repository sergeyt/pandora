package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/sergeyt/pandora/modules/apiutil"
	"github.com/sergeyt/pandora/modules/auth"
	"github.com/sergeyt/pandora/modules/dgraph"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/go-chi/chi"
	authbase "github.com/gocontrib/auth"
	"github.com/gocontrib/pubsub"
)

func dataAPI(r chi.Router) {
	r = r.With(auth.AuthMiddleware)
	r = r.With(dgraph.TransactionMiddleware)

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
		apiutil.SendError(w, err)
		return
	}

	tx := dgraph.RequestTransaction(r)

	resp, err := tx.Query(r.Context(), string(query))
	if err != nil {
		apiutil.SendError(w, err)
		return
	}

	w.Header().Set("Content-Type", apiutil.TypeJSON)
	w.Write(resp.GetJson())
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	resourceType := strings.ToLower(chi.URLParam(r, "type"))
	pg, err := apiutil.ParsePagination(r)
	if err != nil {
		apiutil.SendError(w, err)
		return
	}

	tx := dgraph.RequestTransaction(r)

	data, err := dgraph.ReadList(r.Context(), tx, dgraph.NodeLabel(resourceType), pg)
	if err != nil {
		apiutil.SendError(w, err)
		return
	}

	apiutil.SendJSON(w, data)
}

func meHandler(w http.ResponseWriter, r *http.Request) {
	user := authbase.GetRequestUser(r)
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
	tx := dgraph.RequestTransaction(r)

	data, err := dgraph.ReadNode(r.Context(), tx, id)
	if err != nil {
		apiutil.SendError(w, err)
		return
	}

	apiutil.SendJSON(w, data)
}

func mutateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resourceType := strings.ToLower(chi.URLParam(r, "type"))
	id := chi.URLParam(r, "id")
	nodeLabel := dgraph.NodeLabel(resourceType)
	user := authbase.GetContextUser(ctx)

	var in OrderedJSON
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		apiutil.SendError(w, err)
		return
	}

	isNew := len(id) == 0
	tx := dgraph.RequestTransaction(r)
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
		apiutil.SendError(w, err)
		return
	}

	resp, err := tx.Mutate(ctx, &api.Mutation{
		SetJson: data,
	})
	if err != nil {
		apiutil.SendError(w, err)
		return
	}

	var results []map[string]interface{}

	if isNew {
		results = make([]map[string]interface{}, len(resp.Uids))
		i := 0
		for _, uid := range resp.Uids {
			result, err := dgraph.ReadNode(ctx, tx, uid)
			if err != nil {
				apiutil.SendError(w, err)
				return
			}
			results[i] = result
			i = i + 1
			if len(results) == 1 {
				id = uid
			}
		}
	} else {
		result, err := dgraph.ReadNode(ctx, tx, id)
		if err != nil {
			apiutil.SendError(w, err)
			return
		}
		results = []map[string]interface{}{result}
	}

	err = tx.Commit(ctx)
	if err != nil {
		apiutil.SendError(w, err)
		return
	}

	var out interface{} = results
	if len(results) == 1 {
		out = results[0]
	}

	err = apiutil.SendJSON(w, out)
	if err != nil {
		return
	}

	apiutil.SendEvent(user, &pubsub.Event{
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
	user := authbase.GetContextUser(ctx)

	tx := dgraph.RequestTransaction(r)

	resp, err := tx.Mutate(ctx, &api.Mutation{
		DelNquads: []byte("<" + id + "> * * .\n"),
		CommitNow: true,
	})
	if err != nil {
		apiutil.SendError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)

	apiutil.SendEvent(user, &pubsub.Event{
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
