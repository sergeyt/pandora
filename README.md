# PANDORA

Small app basis composed from the following technologies:

* [Dgraph](https://dgraph.io/) as data store with GraphQL support, write operations using REST
* [tusd](https://tus.io/) as file store baked by Amazon S3 compatible storage like [Minio](https://www.minio.io/)
* [Elasticsearch](https://www.elastic.co/products/elasticsearch) as search engine. Dgraph data is automatically replicated in elasticseach index - not implemented yet :)
* [Kibana](https://www.elastic.co/products/kibana) to visualize Elasticsearch data
* [nats](https://nats.io/) as messaging system with streaming push notifications via [SSE](https://en.wikipedia.org/wiki/Server-sent_events) channel

## Basic Idea

Simple, Flexible, Dynamic, Declarative, Reactive, Realtime Information System :)

## How to run

Just issue `docker-compose up` command in your shell to run all app services

## How to run tests

Run below commands in project directory.

* `go test -coverprofile cover.out` to run tests with coverage output
* `go tool cover -html cover.out` to see HTML coverage report
