# PANDORA

This small box of pandora (aka small app basis) composed from the following technologies:

* [Dgraph](https://dgraph.io/) as data store with GraphQL support, write operations using REST
* [tusd](https://tus.io/) as file store service baked by Amazon S3 compatible storage like [Minio](https://www.minio.io/)
* [ElasticSearch](https://www.elastic.co/products/elasticsearch) as search engine. Dgraph data is automatically replicated in elasticseach index
* [Kibana](https://www.elastic.co/products/kibana) to visualize Elasticsearch data
* [NATS](https://nats.io/) as messaging system with streaming of push notifications (events) via [SSE](https://en.wikipedia.org/wiki/Server-sent_events) channel
* try [FluentD](https://www.fluentd.org/) - for centralized logging, not implemented yet :)

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
1. `imageproxy` - [service](https://willnorris.com/go/imageproxy) with image manipulation ops like resizing
1. `elasticsearch` - search and analitycs engine
1. `kibana` - Elasticsearch dashboard
1. `app` - application API service
1. `caddy` - web server as service gateway

## How to build

To start developing a project, you need to
install git:

    $ sudo apt-get install git

For the `add-apt-repository` to work install package:

    $ sudo apt install software-properties-common

software-properties-common - This software provides an abstraction of the used apt repositories. This allows you to easily manage your distribution and independent software vendors.

Install Docker:

    $ sudo apt update
    $ sudo apt install docker.io docker-compose

Add repository to install Go:

    $ sudo add-apt-repository ppa:longsleep/golang-backports

Install Go:

    $ sudo apt-get update
    $ sudo apt-get install golang-go

Upload project files:

    $ go get github.com/sergeyt/pandora 

Get loads packages called import paths, along with their dependencies. It then installs named packages, such as “go install”.
dep ensure - is the main command and is the only command that changes the state of the disk.

Install dep:

    $ curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

Set project dependencies:

    $ dep ensure

Run project build:

    $ docker-compose up 

Run the script initdata.py (fill the database) execute commands:

    $ apt install python-pip
    $ pip install PyJWT
    $ pip install python-dotenv
    $ pip install Faker

Then execute the script itself:

    $ python initdata.py

## How to run tests

* `go test -coverprofile cover.out`
* `go tool cover -html cover.out`
