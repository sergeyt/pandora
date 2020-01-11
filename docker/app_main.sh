#!/bin/sh

cd /pandora
cat .air.conf
go install
dockerize --wait tcp://dgraph:9080 --wait tcp://minio:9000 --wait tcp://nats:4222 pandora
