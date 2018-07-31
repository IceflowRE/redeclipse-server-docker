#! /bin/bash
# own build script
# load Docker password from ./dockerpassword.txt, which has no new line at the end

echo $(date +%Y-%m-%d=%H:%M:%S)

cat ./dockerpassword.txt | docker login -u=iceflower --password-stdin
mkdir -p ~/.re-docker/sha/re/
mkdir -p ~/.re-docker/sha/baseImg/
mkdir -p ~/.re-docker/sha/dockerfile/

arch="$1"
branches=("stable" "master")

# loops through all branches and update
for branch in "${branches[@]}"; do
	BRANCH="$branch" arch="$arch" ./scripts/update.sh
done
docker logout
