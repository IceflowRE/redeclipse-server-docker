#! /bin/bash
# own build script
# load Docker password from ./dockerpassword.txt, which has no new line at the end

echo $(date +%Y-%m-%d=%H:%M:%S)

docker login -u=iceflower -p="$(cat ./dockerpassword.txt)"
mkdir -p ~/.re-docker/sha/re/
mkdir -p ~/.re-docker/sha/docker/

prefixes=("arm64v8")
branches=("stable" "master")

# loops through all combinations of prefixes and branches and triggers and update check
for prefix in "${prefixes[@]}"; do
    # if prefix is not empty add - after prefix and / for the base image
    if [ $prefix ]; then
        preimage="$prefix/"
        prefix="$prefix-"
    fi
    for branch in "${branches[@]}"; do
        BRANCH="$branch" prefix="$prefix" preimage="$preimage" ./scripts/update.sh
    done
done
