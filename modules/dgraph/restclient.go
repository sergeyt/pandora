package dgraph

import (
	"github.com/gocontrib/rest"
	"github.com/sergeyt/pandora/modules/config"
)

func NewRestClient() *rest.Client {
	return rest.NewClient(rest.Config{
		BaseURL: "http://" + config.DB.Addr,
	})
}
