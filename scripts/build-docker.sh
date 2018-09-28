#!/bin/bash
# build a docker image for the given arch and branch
#
# variable branch, arch and recommit will be passed
# recommit is the latest commit sha for the given branch

repo="iceflower/redeclipse-server"
# later use one dockerfile, but stable has no cmake support atm
#docker build --squash --build-arg BRANCH="$branch" -t "$repo:$arch-$branch" --build-arg RECOMMIT=$recommit -f "Dockerfile_$branch" .
docker build --build-arg BRANCH="$branch" --build-arg RECOMMIT=$recommit -t "$repo:$arch-$branch" -f "Dockerfile_$branch" .
if [ $? -ne 0 ]; then
    exit 1
fi

docker push "$repo:$arch-$branch"
if [ $? -ne 0 ]; then
    exit 1
fi

docker manifest create iceflower/redeclipse-server:$branch iceflower/redeclipse-server:amd64-$branch iceflower/redeclipse-server:arm64v8-$branch
docker manifest annotate iceflower/redeclipse-server:$branch iceflower/redeclipse-server:arm64v8-$branch --variant armv8
docker manifest push --purge iceflower/redeclipse-server:$branch
if [ $? -ne 0 ]; then
    exit 1
fi

exit 0
