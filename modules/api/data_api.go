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
	r = r.With(auth.RequireAPIKey)
	r = r.With(dgraph.TransactionMiddleware)

	r.Post("/api/query", queryHandler)
	r.Get("/api/me", meHandler)
	r.Post("/api/nquads", nquadMutationHandler)
	r.Get("/api/data/{type}/list", listHandler)
	r.Get("/api/data/{type}/{id}", readHandler)

	// mutation api
	r.Post("/api/data/{type}", jsonMutationHandler)
	r.Put("/api/data/{type}/{id}", jsonMutationHandler)
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

	send := func(resp []byte) {
		var m map[string][]interface{}
		err := json.Unmarshal(resp, m)
		if err == nil {
			log.Errorf("json.Unmarshal fail: %v", err)
			apiutil.SendError(w, err)
			return
		}

		empty := true
		for _, v := range m {
			if len(v) > 0 {
				empty = false
				break
			}
		}

		if empty {
			for k, v := range m {
				if len(v) == 0 {
					log.Infof("%s is empty", k)
				}
			}
			http.Error(w, "empty result set", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", apiutil.TypeJSON)
		w.Write(resp)
	}

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

		send(resp.GetJson())
	} else {
		resp, err := tx.Query(r.Context(), query)
		if err != nil {
			log.Errorf("dgraph.Txn.Query fail: %v", err)
			apiutil.SendError(w, err)
			return
		}

		send(resp.GetJson())
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

func jsonMutationHandler(w http.ResponseWriter, r *http.Request) {
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

func nquadMutationHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Errorf("mime.ParseMediaType fail: %v", err)
		apiutil.SendError(w, err)
		return
	}

	var mutation *api.Mutation

	if mediaType == apiutil.TypeJSON {
		var input struct {
			Set    string `json:"set,omitempty"`
			Delete string `json:"delete,omitempty"`
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Errorf("json.Decoder.Decode fail: %v", err)
			apiutil.SendError(w, err)
			return
		}
		if len(input.Set) == 0 && len(input.Delete) == 0 {
			apiutil.SendError(w, fmt.Errorf("invalid input. please specify set or delete mutations"))
		}

		mutation = &api.Mutation{}
		if len(input.Set) > 0 {
			mutation.SetNquads = []byte(input.Set)
		}
		if len(input.Delete) > 0 {
			mutation.DelNquads = []byte(input.Delete)
		}
	} else {
		nquads, err := ioutil.ReadAll(r.Body)
		if err != nil {
			apiutil.SendError(w, err, http.StatusInternalServerError)
			return
		}

		mutation = &api.Mutation{
			SetNquads: nquads,
		}
	}

	ctx := r.Context()
	user := authbase.GetContextUser(ctx)
	// TODO add metadata props created_by, modified_by
	tx := dgraph.RequestTransaction(r)

	mutation.CommitNow = true
	resp, err := tx.Mutate(ctx, mutation)
	if err != nil {
		log.Errorf("dgraph.Txn.Mutate fail: %v", err)
		apiutil.SendError(w, err, http.StatusInternalServerError)
		return
	}

	err = apiutil.SendJSON(w, resp)
	if err != nil {
		return
	}

	now := time.Now()

	for _, v := range resp.Uids {
		apiutil.SendEvent(user, &pubsub.Event{
			Action:     r.Method,
			Method:     r.Method,
			URL:        r.URL.String(),
			ResourceID: v,
			CreatedBy:  user.GetID(),
			CreatedAt:  now,
		})
	}
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
