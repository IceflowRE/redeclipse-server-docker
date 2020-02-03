#!/bin/bash -e

docker-compose stop -t 600 master stable stable-re2
#docker-compose down
rsd-updater /home/iceflower/.re-updater/ --config
docker-compose -p re-server up -d master
docker-compose -p re-server up -d stable
docker-compose -p re-server up -d stable-re2
