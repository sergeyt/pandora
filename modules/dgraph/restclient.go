package dgraph

import (
	"net/http"
	"os"

	"github.com/gocontrib/rest"
)

func NewRestClient() *rest.Client {
	return rest.NewClient(rest.Config{
		// TODO get from env
		BaseURL: "http://dgraph:8080",
		Verbose: true,
		Headers: func(h http.Header) {
			token := os.Getenv("DGRAPH_TOKEN")
			if token != "" {
				h.Set("X-Dgraph-AuthToken", token)
			}
		},
	})
}
