package main

import (
	"fmt"
	stdlog "log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/gocontrib/request"
	"github.com/tus/tusd"
	"github.com/tus/tusd/filestore"
)

func uploadAPI(mux chi.Router) {
	dir := "./data/uploads"
	os.MkdirAll(dir, 0700)

	// Create a new FileStore instance which is responsible for
	// storing the uploaded file on disk in the specified directory.
	// If you want to save them on a different medium, for example
	// a remote FTP server, you can implement your own storage backend
	// by implementing the tusd.DataStore interface.
	store := filestore.FileStore{
		Path: dir,
	}

	logger := stdlog.New(os.Stdout, "[tusd] ", stdlog.LstdFlags)

	config := tusd.Config{
		BasePath:  "/api/uploads/",
		DataStore: store,
		Logger:    logger,
	}

	h, err := tusd.NewUnroutedHandler(config)
	if err != nil {
		panic(fmt.Sprintf("unable to create upload handler: %s", err))
	}

	mux = mux.With(authMiddleware)
	mux.Post("/api/uploads", tusdHandler(h, h.PostFile))
	mux.Head("/api/uploads/:id", tusdHandler(h, h.HeadFile))
	mux.Get("/api/uploads/:id", tusdHandler(h, h.GetFile))
	mux.Patch("/api/uploads/:id", tusdHandler(h, h.PatchFile))

	// Only attach the DELETE handler if the Terminate() method is provided
	if _, ok := config.DataStore.(tusd.TerminaterDataStore); ok {
		mux.Delete("/api/uploads/:id", tusdHandler(h, h.DelFile))
	}
}

func tusdHandler(handler *tusd.UnroutedHandler, fn http.HandlerFunc) http.HandlerFunc {
	h := handler.Middleware(fn)
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO fix r.Host on top middleware
		// tus is not proxy friendly since it uses r.Host to make absolute URL to uploaded file
		// as workaround we replace r.Host value with possible X-Forwarded-Host value
		host := request.GetHost(r)
		if host != r.Host {
			r.Host = host
		}
		h.ServeHTTP(w, r)
	}
}
