#!/usr/bin/env bash

set -e

git checkout Gopkg.lock
git pull
dep ensure
go install

SERVICES=`docker-compose ps --services`
while read -r SERVICE; do
    if [ "$SERVICE" -ne "caddy" ]
    then
        docker-compose stop $SERVICE
    fi
done <<< "$SERVICES"

docker-compse up -d
docker-compose ps
