package elasticsearch

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gocontrib/esclient"
	"github.com/sergeyt/pandora/modules/apiutil"
	"github.com/sergeyt/pandora/modules/auth"
)

// SearchAPI is proxy os ElasticSearch query
// TODO just use http.Proxy or proxy with caddy
func SearchAPI(r chi.Router) {
	r = r.With(auth.Middleware)

	r.Get("/api/search/:idx", func(w http.ResponseWriter, r *http.Request) {
		idxName := chi.URLParam(r, "idx")

		var sr esclient.SearchRequest
		err := json.NewDecoder(r.Body).Decode(&sr)
		if err != nil {
			apiutil.SendError(w, err, http.StatusBadRequest)
			return
		}

		c := makeClient()
		result, err := c.Search(idxName, &sr, nil)
		if err != nil {
			apiutil.SendError(w, err)
			return
		}

		apiutil.SendJSON(w, result)
	})
}
