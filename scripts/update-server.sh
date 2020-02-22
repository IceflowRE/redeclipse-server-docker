#!/bin/bash -e

cd /home/iceflower/redeclipse-server-data/
docker-compose stop -t 600
docker-compose pull
docker-compose up -d
