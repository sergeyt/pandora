package main

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/go-chi/chi"
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
	// TODO read all values of given node
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

	sendJSON(w, resp)

	// TODO notify about change
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	tx := transaction(r)
	id := chi.URLParam(r, "id")

	resp, err := tx.Mutate(r.Context(), &api.Mutation{
		DelNquads: []byte("<" + id + "> * * .\n"),
		CommitNow: true,
	})
	if err != nil {
		sendError(w, err)
		return
	}

	sendJSON(w, resp)

	// TODO notify about change
}
