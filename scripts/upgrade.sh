#!/usr/bin/env bash

docker-compose down

git checkout Gopkg.lock
git pull
dep ensure
go install

if [ $REBUILD = "yes" ]; then
    docker-compose build --no-cache
fi

docker-compose up -d
docker-compose ps
