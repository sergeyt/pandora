package elasticsearch

import (
	"github.com/gocontrib/esclient"
	"github.com/gocontrib/log"
	"github.com/sergeyt/pandora/modules/config"
)

func makeClient() *esclient.Client {
	return esclient.NewClient(esclient.Config{
		URL: config.ElasticSearch.URL,
	})
}

func EnsureIndex() {
	c := makeClient()
	name := config.ElasticSearch.IndexName

	if !c.IndexExists(name) {
		err := c.CreateIndex(name)
		if err != nil {
			log.Errorf("elasticseach: cannot create index: %v", err)
		} else {
			log.Info("elasticseach: created index %s", name)
		}
	}
}
