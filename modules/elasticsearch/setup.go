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

func EnsureIndex() {
	c := makeClient()
	if !c.IndexExists(config.ElasticSearch.IndexName) {
		c.CreateIndex(config.ElasticSearch.IndexName)
	}
}
