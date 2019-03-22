#!/usr/bin/env bash

set -e

git checkout Gopkg.lock
git pull
dep ensure
go install

SERVICES=`docker-compose ps --services`
while read -r SERVICE; do
    echo "stopping ${SERVICE}"
    if [ "${SERVICE}" != "caddy" ]
    then
        docker-compose stop $SERVICE
    fi
done <<< "$SERVICES"

docker-compose up -d
docker-compose restart caddy
docker-compose ps
