#!/usr/bin/env bash

git checkout Gopkg.lock
git stash
git pull
git stash apply
go install
dep ensure

docker-compose down

if [ $REBUILD = "yes" ]; then
    docker-compose build --no-cache
fi

docker-compose up -d
docker-compose ps
