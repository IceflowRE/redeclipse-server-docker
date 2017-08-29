#! /bin/bash

# variable BRANCH is defined inside the travis.yml env or will be passed
# variable prefix for e.g. architectures will be passed

# load latest branch commit sha and get latest from remote
savedSha="$(cat ./sha/re/$BRANCH.sha)"
sha="$(git ls-remote --heads https://github.com/red-eclipse/base.git $BRANCH | awk '{ print $1 }')"
if [ "$sha" == "" ]; then
    echo "Cant get latest git commit sha."
    exit 1
fi

# get alpine docker image, since its only some mb its ok
docker pull alpine
# load saved base image sha and get latest image sha
savedBaseImgSha="$(cat ./sha/docker/alpine.sha)"
baseImgSha="$(docker image ls --digests --format '{{.Digest}}' alpine)"
if [ "$baseImgSha" == "" ]; then
    echo "Cant get latest docker sha."
    exit 1
fi

echo "git: $savedSha | $sha"
echo "img: $savedBaseImgSha | $baseImgSha"

# update only if a new commit exists or the base image was updated
if [ "$savedSha" != "$sha" ] || [ "$savedBaseImgSha" != "$baseImgSha" ]; then
    echo "Build $BRANCH"
    branch="$BRANCH" prefix=$prefix ./scripts/build-docker.sh
    if [ $? -ne 0 ]; then # dont save sha if something failed
        exit 1
    fi
    # save latest shas
    echo "Save sha $sha as $BRANCH"
    echo "$sha" > "./sha/re/$BRANCH.sha"
    echo "Save base image sha $baseImgSha as alpine"
    echo "$baseImgSha" > "./sha/docker/alpine.sha"
else
    echo "Skip $BRANCH - $sha | alpine - $baseImgSha"
fi
exit 0
