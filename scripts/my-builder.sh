#! /bin/bash
# own build script
# load Docker password from ./dockerpassword.txt, which has no new line at the end

docker login -u=iceflower -p="$(cat ./dockerpassword.txt)"
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
