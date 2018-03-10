# PANDORA

## How to run tests

Run below commands in project directory.

* `go test -coverprofile cover.out` to run tests with coverage output
* `go tool cover -html cover.out` to see HTML coverage report

Extensible dynamic content management system build on top of communication concepts:

* users
* messages
* channels

### Features ###

* built-in real-time messaging
* Git like communication model

### Architecture ###

* caddyserver - web server as gateway
* dgraph - fast graph database
* nats - fast pubsub service
* search engine - elastic, bleve?
* thin golang web server
* server sent events to get updates of
* subscribe on messages/channels

### Week 1 ###

* learn dgraph
* create schema for micro chat
* no authentication
* SSE of changes in dgraph, probably hacking dgraph
* try to build simple micro chat

### TODO ### 

* simple login - https://gist.github.com/mschoebel/9398202
