#! /bin/bash
# checks for an update and triggers a build then
# as a check local saved checksums will be compared to the latest checksums
# if an update is available a build is triggered and everything wents well, the new checksum will be saved
#
# variable BRANCH is defined inside the travis.yml env or will be passed
# variable prefix for e.g. architectures will be passed

# load latest branch commit sha and get latest from remote
savedSha="$(cat ~/.re-docker/sha/re/$prefix$BRANCH.sha)"
sha="$(git ls-remote --heads https://github.com/red-eclipse/base.git $BRANCH | awk '{ print $1 }')"
if [ "$sha" == "" ]; then
    echo "Cant get latest git commit sha."
    exit 1
fi

# get alpine docker image, since its only some mb its ok
docker pull "$preimage"alpine

# load saved base image sha and get latest image sha
savedBaseImgSha="$(cat ~/.re-docker/sha/docker/$prefix$BRANCH-alpine.sha)"
baseImgSha="$(docker image ls --digests --format '{{.Digest}}' alpine)"
if [ "$baseImgSha" == "" ]; then
    echo "Cant get latest docker sha."
    exit 1
fi

echo "git: $savedSha | $sha"
echo "img: $savedBaseImgSha | $baseImgSha"

# update only if a new commit exists or the base image was updated
if [ "$savedSha" != "$sha" ] || [ "$savedBaseImgSha" != "$baseImgSha" ]; then
    echo "Build $prefix$BRANCH"
    branch="$BRANCH" prefix="$prefix" preimage="$preimage" ./scripts/build-docker.sh
    if [ $? -ne 0 ]; then # dont save sha if something failed
        exit 1
    fi
    # save latest shas
    echo "Save sha $sha as $prefix$BRANCH"
    echo "$sha" > ~/.re-docker/sha/re/"$prefix$BRANCH".sha
    echo "Save base image sha $baseImgSha as alpine"
    echo "$baseImgSha" > ~/.re-docker/sha/docker/"$prefix$BRANCH"-alpine.sha
else
    echo "Skip $prefix$BRANCH"
fi
exit 0
