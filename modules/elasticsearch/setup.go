package elasticsearch

import (
	"github.com/gocontrib/esclient"
	"github.com/sergeyt/pandora/modules/config"
)

func makeClient() *esclient.Client {
	return esclient.NewClient(esclient.Config{
		URL: config.ElasticSearch.URL,
	})
}
