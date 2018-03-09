package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func makeAPIHandler() http.Handler {
	r := chi.NewRouter()
	r.Get("/query", queryHandler)
	r.Post("/{type}", createHandler)
	r.Put("/{type}/{id}", updateHandler)
	r.Delete("/{type}/{id}", deleteHandler)
	return r
}

func queryHandler(w http.ResponseWriter, r *http.Request) {

}

func createHandler(w http.ResponseWriter, r *http.Request) {

}

func updateHandler(w http.ResponseWriter, r *http.Request) {

}

func deleteHandler(w http.ResponseWriter, r *http.Request) {

}
