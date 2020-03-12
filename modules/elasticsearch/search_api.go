package elasticsearch

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gocontrib/esclient"
	"github.com/sergeyt/pandora/modules/auth"
	"github.com/sergeyt/pandora/modules/send"
)

// SearchAPI is proxy os ElasticSearch query
// TODO just use http.Proxy or proxy with caddy
func SearchAPI(r chi.Router) {
	r = r.With(auth.Middleware)

	do := func(w http.ResponseWriter, r *http.Request, sr *esclient.SearchRequest) {
		idxName := chi.URLParam(r, "idx")

		c := makeClient()
		result, err := c.Search(idxName, sr, nil)
		if err != nil {
			send.Error(w, err)
			return
		}

		send.JSON(w, result)
	}

	r.Get("/api/search/:idx", func(w http.ResponseWriter, r *http.Request) {
		sr := esclient.SearchRequest{
			Query: map[string]interface{}{
				"match": map[string]string{
					"text": r.URL.Query().Get("query"),
				},
			},
		}
		do(w, r, &sr)
	})

	r.Post("/api/search/:idx", func(w http.ResponseWriter, r *http.Request) {
		var sr esclient.SearchRequest
		err := json.NewDecoder(r.Body).Decode(&sr)
		if err != nil {
			send.Error(w, err, http.StatusBadRequest)
			return
		}

		do(w, r, &sr)
	})
}
