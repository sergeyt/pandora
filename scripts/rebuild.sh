#!/usr/bin/env bash

set -e

git checkout Gopkg.lock
git pull

docker-compose down

dep ensure
go install

docker-compose build --no-cache
docker-compose up -d
# TODO wait while all services is up
docker-compose ps
