#!/bin/bash
# build a docker image for the given prefix and branch
#
# variable branch and prefix will be passed

repo="iceflower/redeclipse-server"
# later use one dockerfile, but stable has no cmake support atm
#docker build --squash --build-arg BRANCH="$branch" -t "$repo:$prefix$branch" -f "Dockerfile_$branch" .
docker build --build-arg PREIMAGE="$preimage" --build-arg BRANCH="$branch" -t "$repo:$prefix$branch" -f "Dockerfile_$branch" .
if [ $? -ne 0 ]; then
    exit 1
fi

docker push "$repo:$prefix$branch"
if [ $? -ne 0 ]; then
    exit 1
fi
exit 0
