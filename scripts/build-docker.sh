#!/bin/bash

repo="iceflower/red-eclipse_devel_server_test"
branch="$1"  # master or stable, given as argument
# later use one dockerfile, but stable has no cmake support atm
#docker build --squash -t "$repo:$branch" -f "Dockerfile_$branch" .
docker build --build-arg BRANCH="$branch" -t "$repo:$branch" -f "Dockerfile_$branch" .
if [ $? -ne 0 ]; then
    exit 1
fi
docker push "$repo:$branch"
if [ $? -ne 0 ]; then
    exit 1
fi
exit 0
