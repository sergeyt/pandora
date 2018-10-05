package elasticsearch

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/gocontrib/esclient"
	"github.com/sergeyt/pandora/modules/apiutil"
	"github.com/sergeyt/pandora/modules/config"
)

// TODO just use http.Proxy or proxy with caddy
func SearchAPI(r chi.Router) {
	r.Get("/api/search", func(w http.ResponseWriter, r *http.Request) {
		var sr esclient.SearchRequest
		err := json.NewDecoder(r.Body).Decode(&sr)
		if err != nil {
			apiutil.SendError(w, err, http.StatusBadRequest)
			return
		}

		c := makeClient()
		result, err := c.Search(config.ElasticSearch.IndexName, &sr, nil)
		if err != nil {
			apiutil.SendError(w, err)
			return
		}

		apiutil.SendJSON(w, result)
	})
}
