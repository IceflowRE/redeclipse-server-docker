#!/bin/bash

REPO="iceflower/red-eclipse_devel_server_test"
BRANCH="$1"  # master or stable, given as argument
# later use one dockerfile, but stable has no cmake support atm
#docker build --squash -t "$REPO:$BRANCH" -f "Dockerfile_$BRANCH" .
docker build -t "$REPO:$BRANCH" -f "Dockerfile_$BRANCH" .
docker push "$REPO:$BRANCH"
exit 0
