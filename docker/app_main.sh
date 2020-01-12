#!/bin/sh

cd /pandora
go install
dockerize --wait tcp://dgraph:9080 --wait tcp://minio:9000 --wait tcp://nats:4222 /usr/local/bin/air
