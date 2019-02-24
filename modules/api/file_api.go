package api

import (
	"fmt"
	"io"
	"mime"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sergeyt/pandora/modules/apiutil"
	"github.com/sergeyt/pandora/modules/auth"
	log "github.com/sirupsen/logrus"
)

func fileAPI(r chi.Router) {
	r = r.With(auth.Middleware)

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

		store := NewFileStore()
		h(fsopContext{store, path}, w, r)
	}
}

type fsopContext struct {
	store ObjectStore
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
	// TODO mime type filter

	ct := r.Header.Get("Content-Type")
	mt, _, err := mime.ParseMediaType(ct)
	if err != nil {
		log.Errorf("mime.ParseMediaType fail: %v", err)
		apiutil.SendError(w, err)
		return
	}
	if mt == "multipart/form-data" || mt == "multipart/mixed" {
		mr, err := r.MultipartReader()
		if err != nil {
			log.Errorf("http.Request.MultipartReader fail: %v", err)
			apiutil.SendError(w, err)
			return
		}
		results := make(map[string]interface{})
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Errorf("multipart.Reader.NextPart fail: %v", err)
				apiutil.SendError(w, err)
				return
			}

			path := c.path
			if len(path) > 0 {
				path = path + "/" + p.FileName()
			}

			result, err := c.store.Upload(r.Context(), path, p)
			if err != nil {
				log.Errorf("FileStore.Upload fail: %v", err)
				apiutil.SendError(w, err)
				return
			}

			results[path] = result
		}

		if len(results) > 0 {
			apiutil.SendJSON(w, results)
		}

		return
	}

	// TODO generate filename if path is not defined

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
	result, err := c.store.Delete(r.Context(), c.path)
	if err != nil {
		log.Errorf("FileStore.Delete fail: %v", err)
		apiutil.SendError(w, err)
		return
	}

	if result != nil {
		apiutil.SendJSON(w, result)
	}
}
