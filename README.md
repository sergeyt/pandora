# PANDORA

This small box of pandora (IOW small app basis) composed from the following technologies:

* [Dgraph](https://dgraph.io/) as data store with GraphQL support, write operations using REST
* [tusd](https://tus.io/) as file store service baked by Amazon S3 compatible storage like [Minio](https://www.minio.io/)
* [Elasticsearch](https://www.elastic.co/products/elasticsearch) as search engine. Dgraph data is automatically replicated in elasticseach index
* [Kibana](https://www.elastic.co/products/kibana) to visualize Elasticsearch data
* [nats](https://nats.io/) as messaging system with streaming of push notifications (events) via [SSE](https://en.wikipedia.org/wiki/Server-sent_events) channel
* try [fluentd](https://www.fluentd.org/) - for centralized logging, not implemented yet :)

## Basic Idea

I'd like to have simple, flexible, dynamic, declarative, reactive, realtime information system :)

## How to run

`docker-compose up` runs all app services:

1. `zero` - Dgraph cluster manager
1. `dgraph` - Dgraph data manager hosts predicates & indexes
1. `ratel` - serves the UI to run queries, mutations & altering schema
1. `nats` - plays as message bus
1. `pubsub` - event streaming service based on [SSE](https://en.wikipedia.org/wiki/Server-sent_events) protocol
1. `minio` - Amazon S3 compatible file store
1. `tusd` - service with Open Protocol for Resumable File Uploads
1. `elasticsearch` - search and analitycs engine
1. `kibana` - Elasticsearch dashboard
1. `app` - application API service
1. `caddy` - web server as service gateway

## How to build

* `dep ensure`
* `go install`

## How to run tests

* `go test -coverprofile cover.out`
* `go tool cover -html cover.out`
