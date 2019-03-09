package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/dgraph-io/dgo/protos/api"
	"github.com/sergeyt/pandora/modules/apiutil"
	"github.com/sergeyt/pandora/modules/auth"
	"github.com/sergeyt/pandora/modules/dgraph"
	"github.com/sergeyt/pandora/modules/utils"
	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	authbase "github.com/gocontrib/auth"
	"github.com/gocontrib/pubsub"
)

func dataAPI(r chi.Router) {
	r = r.With(auth.Middleware)
	r = r.With(dgraph.TransactionMiddleware)

	r.Post("/api/query", queryHandler)
	r.Get("/api/me", meHandler)
	r.Post("/api/nquads", setNquads)
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("ioutil.ReadAll fail: %v", err)
		apiutil.SendError(w, err)
		return
	}

	query := string(body)
	// TODO use regexp for better matching
	hasVars := strings.Contains(query, "$")

	tx := dgraph.RequestTransaction(r)

	if hasVars {
		vars := make(map[string]string)
		for k, v := range r.URL.Query() {
			if strings.HasPrefix(k, "$") {
				vars[k] = v[0]
			}
		}

		resp, err := tx.QueryWithVars(r.Context(), query, vars)
		if err != nil {
			log.Errorf("dgraph.Txn.QueryWithVars fail: %v", err)
			apiutil.SendError(w, err)
			return
		}

		w.Header().Set("Content-Type", apiutil.TypeJSON)
		w.Write(resp.GetJson())
	} else {
		resp, err := tx.Query(r.Context(), query)
		if err != nil {
			log.Errorf("dgraph.Txn.Query fail: %v", err)
			apiutil.SendError(w, err)
			return
		}

		w.Header().Set("Content-Type", apiutil.TypeJSON)
		w.Write(resp.GetJson())
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	resourceType := strings.ToLower(chi.URLParam(r, "type"))
	pg, err := apiutil.ParsePagination(r)
	if err != nil {
		log.Errorf("apiutil.ParsePagination fail: %v", err)
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
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Errorf("mime.ParseMediaType fail: %v", err)
		apiutil.SendError(w, err)
		return
	}

	ctx := r.Context()
	resourceType := strings.ToLower(chi.URLParam(r, "type"))
	id := chi.URLParam(r, "id")
	nodeLabel := dgraph.NodeLabel(resourceType)
	user := authbase.GetContextUser(ctx)

	if mediaType == "application/json" {

		var in utils.OrderedJSON
		err = json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			log.Errorf("bad JSON. json.Decoder.Decode fail: %v", err)
			apiutil.SendError(w, err)
			return
		}

		tx := dgraph.RequestTransaction(r)

		results, err := dgraph.Mutate(ctx, tx, dgraph.Mutation{
			Input:     in,
			NodeLabel: nodeLabel,
			ID:        id,
			By:        user.GetID(),
		})
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
	} else {
		apiutil.SendError(w, fmt.Errorf("unsupported media type: %s", contentType), http.StatusUnsupportedMediaType)
	}
}

func setNquads(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Errorf("mime.ParseMediaType fail: %v", err)
		apiutil.SendError(w, err)
		return
	}

	if mediaType != "application/n-quads" {
		apiutil.SendError(w, fmt.Errorf("unsupported media type: %s", contentType), http.StatusUnsupportedMediaType)
		return
	}

	nquads, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apiutil.SendError(w, err, http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	// TODO add metadata props created_by, modified_by
	tx := dgraph.RequestTransaction(r)

	resp, err := tx.Mutate(ctx, &api.Mutation{
		SetNquads: nquads,
		CommitNow: true,
	})
	if err != nil {
		log.Errorf("dgraph.Txn.Mutate fail: %v", err)
		apiutil.SendError(w, err, http.StatusInternalServerError)
		return
	}

	// TODO raise pubsub event
	apiutil.SendJSON(w, resp)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	resourceType := strings.ToLower(chi.URLParam(r, "type"))
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	user := authbase.GetContextUser(ctx)

	tx := dgraph.RequestTransaction(r)

	var fileNode map[string]interface{}

	if resourceType == "file" {
		node, err := dgraph.ReadNode(ctx, tx, id)
		if err != nil {
			apiutil.SendError(w, err)
			return
		}
		fileNode = node
	}

	resp, err := dgraph.DeleteNode(ctx, tx, id)
	if err != nil {
		apiutil.SendError(w, err)
		return
	}

	if fileNode != nil {
		go deleteFileObject(fileNode)
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
