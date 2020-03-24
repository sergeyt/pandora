package api

import (
	"net/url"
	"time"

	health "github.com/InVisionApp/go-health/v2"
	"github.com/InVisionApp/go-health/v2/checkers"
	"github.com/InVisionApp/go-health/v2/handlers"
	"github.com/go-chi/chi"
	"github.com/sergeyt/pandora/modules/config"
	"github.com/sergeyt/pandora/modules/dgraph"
	log "github.com/sirupsen/logrus"
)

func healthAPI(r chi.Router) {
	h := health.New()

	dgraphURL, err := url.Parse(config.DGraph.RpcURL)
	if err != nil {
		log.Fatalf("invalid dgraph URL: %v", err)
	}

	dgraph, err := checkers.NewReachableChecker(&checkers.ReachableConfig{
		URL:    dgraphURL,
		Dialer: dgraph.Dial,
	})
	if err != nil {
		log.Fatalf("NewReachableChecker fail for dgraph: %v", err)
	}
	log.Infof("dgraph URL: %s", dgraphURL)
	log.Infof("dgraph port: %s", dgraphURL.Port())

	natsURL, err := url.Parse(config.NatsURL)
	if err != nil {
		log.Fatalf("invalid NATS URL: %v", err)
	}

	nats, err := checkers.NewReachableChecker(&checkers.ReachableConfig{
		URL: natsURL,
	})
	if err != nil {
		log.Fatalf("NewReachableChecker fail for nats: %v", err)
	}

	inerval := time.Duration(10) * time.Second

	h.AddChecks([]*health.Config{
		{
			Name:     "dgraph",
			Checker:  dgraph,
			Interval: inerval,
			Fatal:    true,
		},
		{
			Name:     "nats",
			Checker:  nats,
			Interval: inerval,
			Fatal:    false,
		},
	})

	if err := h.Start(); err != nil {
		log.Fatalf("unable to start healthcheck: %v", err)
	}

	r.Head("/api/health", handlers.NewBasicHandlerFunc(h))
	r.Get("/api/health", handlers.NewJSONHandlerFunc(h, nil))
}
