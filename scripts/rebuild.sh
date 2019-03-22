#!/usr/bin/env bash

set -e

git checkout Gopkg.lock
git pull
dep ensure
go install

docker-compose down
docker-compose build --no-cache
docker-compse up -d
# TODO wait while all services is up
docker-compose ps
