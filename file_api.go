package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sergeyt/pandora/modules/apiutil"
	"github.com/sergeyt/pandora/modules/auth"
	log "github.com/sirupsen/logrus"
)

func fileAPI(r chi.Router) {
	r = r.With(auth.AuthMiddleware)

	r.Get("/api/file/*", asHTTPHandler(downloadFile))
	r.Post("/api/file/*", asHTTPHandler(uploadFile))
	r.Put("/api/file/*", asHTTPHandler(uploadFile))
	r.Delete("/api/file/*", asHTTPHandler(deleteFile))
}

func asHTTPHandler(h fileHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := chi.URLParam(r, "*")
		err := validatePath(r, path)
		if err != nil {
			log.Errorf("path is not valid: %v", err)
			apiutil.SendError(w, err)
			return
		}

		store, err := makeFileStore()
		if err != nil {
			log.Errorf("makeFileStore fail: %v", err)
			apiutil.SendError(w, err)
			return
		}

		h(fsopContext{store, path}, w, r)
	}
}

type fsopContext struct {
	store FileStore
	path  string
}

type fileHandler func(c fsopContext, w http.ResponseWriter, r *http.Request)

func validatePath(r *http.Request, path string) error {
	if len(path) == 0 {
		if r.Method == "GET" || r.Method == "DELETE" {
			return fmt.Errorf("file path is not defined")
		}
	}
	return nil
}

func downloadFile(c fsopContext, w http.ResponseWriter, r *http.Request) {
	err := c.store.Download(r.Context(), c.path, w)
	if err != nil {
		log.Errorf("FileStore.Download fail: %v", err)
		apiutil.SendError(w, err)
		return
	}
}

func uploadFile(c fsopContext, w http.ResponseWriter, r *http.Request) {
	// TODO generate filename if path is not defined
	// TODO mime type filter

	result, err := c.store.Upload(r.Context(), c.path, r.Body)
	if err != nil {
		log.Errorf("FileStore.Upload fail: %v", err)
		apiutil.SendError(w, err)
		return
	}

	if result != nil {
		apiutil.SendJSON(w, result)
	}
}

func deleteFile(c fsopContext, w http.ResponseWriter, r *http.Request) {
	err := c.store.Delete(r.Context(), c.path)
	if err != nil {
		log.Errorf("FileStore.Delete fail: %v", err)
		apiutil.SendError(w, err)
		return
	}
}
