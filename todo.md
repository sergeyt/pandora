# todo

* [x] basic auth
* [x] server sent events (SSE) about data changes
* [x] label nodes
* [x] auth: match only user nodes
* [x] automatically set created, modified predicates
* [ ] allow using JSON schema to validate inputs in mutations
* [ ] method to modify graph schema
* [ ] select good logger
* [ ] complete dockerization
* [ ] declarative constraints (e.g. uniq)
* [ ] [RBAC](https://en.wikipedia.org/wiki/Role-based_access_control)

## apps

* [ ] small chat with terminal like web shell, e.g. with command to change current channel
* [ ] built-in todo list channel

## cleanups

* [x] (minor) consider moving SSE handler to github.com/gocontrib/pubsub
* [x] SSE as separate upstream process (psubd)
* [x] tusd as separate upstream process to handle file uploads/downloads

## easy todo

* [ ] stream push notifications over web sockets
