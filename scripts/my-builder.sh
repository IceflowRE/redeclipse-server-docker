#! /bin/bash
# own build script

docker login -u=iceflower -p="$DOCKER_PASSWORD"
mkdir -p ~/.re-docker/sha/re/
mkdir -p ~/.re-docker/sha/docker/

prefixes=("arm64v8-")
branches=("stable" "master")

# loops through all combinations of prefixes and branches and triggers and update check
for prefix in "${prefixes[@]}"; do
    for branch in "${branches[@]}"; do
        BRANCH="$branch" prefix="$prefix" ./scripts/update.sh
    done
done
