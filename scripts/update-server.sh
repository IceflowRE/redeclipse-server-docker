#!/bin/bash -e

docker-compose stop -t 600 master stable
#docker-compose down
./scripts/my-builder.sh
docker-compose -p re-server up -d master
docker-compose -p re-server up -d stable
