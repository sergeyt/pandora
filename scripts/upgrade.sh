#!/usr/bin/env bash

set -e

SERVICES=`docker-compose ps --services`

stop_services() {
    while read -r SERVICE; do
        echo "stopping ${SERVICE}"
        if [ "${SERVICE}" != "caddy" ]
        then
            docker-compose stop $SERVICE
        fi
    done <<< "$SERVICES"
}

git checkout Gopkg.lock
git pull

stop_services

dep ensure
go install

docker-compose up -d
docker-compose restart caddy
docker-compose ps
