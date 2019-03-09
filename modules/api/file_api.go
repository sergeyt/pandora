package api

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	Auth "github.com/gocontrib/auth"
	"github.com/gocontrib/pubsub"
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

	ctx := r.Context()
	user := Auth.GetContextUser(ctx)
	// FIXME determine by content type
	resourceType := "file"

	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Errorf("mime.ParseMediaType fail: %v", err)
		apiutil.SendError(w, err)
		return
	}
	if mediaType == "multipart/form-data" || mediaType == "multipart/mixed" {
		mr, err := r.MultipartReader()
		if err != nil {
			log.Errorf("http.Request.MultipartReader fail: %v", err)
			apiutil.SendError(w, err)
			return
		}
		results := make(map[string]map[string]interface{})
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

			result, err := c.store.Upload(ctx, path, mediaType, p)
			if err != nil {
				log.Errorf("FileStore.Upload fail: %v", err)
				apiutil.SendError(w, err)
				return
			}

			results[path] = result
		}

		if len(results) > 0 {
			err = apiutil.SendJSON(w, results)
			if err != nil {
				return
			}

			id := ""
			if len(results) == 1 {
				for _, v := range results {
					id = getUID(v)
					break
				}
			}

			// FIXME send multiple events
			apiutil.SendEvent(user, &pubsub.Event{
				Action:       r.Method,
				Method:       r.Method,
				URL:          r.URL.String(),
				ResourceID:   id,
				ResourceType: resourceType,
				CreatedBy:    user.GetID(),
				CreatedAt:    time.Now(),
				Result:       results,
			})
		}

		return
	}

	// TODO generate filename if path is not defined

	result, err := c.store.Upload(r.Context(), c.path, mediaType, r.Body)
	if err != nil {
		log.Errorf("FileStore.Upload fail: %v", err)
		apiutil.SendError(w, err)
		return
	}

	if result != nil {
		err = apiutil.SendJSON(w, result)
		if err != nil {
			return
		}

		id := getUID(result)

		apiutil.SendEvent(user, &pubsub.Event{
			Action:       r.Method,
			Method:       r.Method,
			URL:          r.URL.String(),
			ResourceID:   id,
			ResourceType: resourceType,
			CreatedBy:    user.GetID(),
			CreatedAt:    time.Now(),
			Result:       result,
		})
	}
}

func deleteFile(c fsopContext, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := Auth.GetContextUser(ctx)
	// FIXME determine by content type
	resourceType := "file"

	id, result, err := c.store.Delete(r.Context(), c.path)
	if err != nil {
		log.Errorf("FileStore.Delete fail: %v", err)
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
		Result:       result,
	})
}

func deleteFileObject(fileNode map[string]interface{}) {
	ctx := context.Background()
	path := fileNode["path"].(string)
	store := NewFileStore()
	store.DeleteObject(ctx, path)
}

func getUID(result map[string]interface{}) string {
	v, ok := result["uid"]
	if ok {
		s, ok := v.(string)
		if ok {
			return s
		}
	}
	return ""
}
