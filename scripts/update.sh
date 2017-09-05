#! /bin/bash
# checks for an update and triggers a build then
# as a check local saved checksums will be compared to the latest checksums
# if an update is available a build is triggered and everything wents well, the new checksum will be saved
# checks are: Dockerfile, base image or redeclipse git changes
#
# variable BRANCH is defined inside the travis.yml env or will be passed
# variable prefix for e.g. architectures will be passed

# load latest branch commit sha and get latest from remote
savedReSha="$(cat ~/.re-docker/sha/re/$prefix$BRANCH.sha)"
reSha="$(git ls-remote --heads https://github.com/red-eclipse/base.git $BRANCH | awk '{ print $1 }')"
if [ "$reSha" == "" ]; then
    echo "Cant get latest git commit sha."
    exit 1
fi

# get alpine docker image, since its only some mb its ok
docker pull "$preimage"alpine

# load saved base image sha and get latest image sha
savedBaseImgSha="$(cat ~/.re-docker/sha/baseImg/$prefix$BRANCH-alpine.sha)"
baseImgSha="$(docker image ls --digests --format '{{.Digest}}' $preimage\alpine)"
if [ "$baseImgSha" == "" ]; then
    echo "Cant get latest base image sha."
    exit 1
fi

# load dockerfile sha and get latest dockerfile sha
savedDockerSha="$(cat ~/.re-docker/sha/dockerfile/Dockerfile_$BRANCH.sha)"
dockerSha="$(sha256sum Dockerfile_$BRANCH | awk '{ print $1 }')"
if [ "$dockerSha" == "" ]; then
    echo "Cant get latest dockerfile sha."
    exit 1
fi

echo "git: $savedReSha | $reSha"
echo "img: $savedBaseImgSha | $baseImgSha"
echo "dockerfile: $savedDockerSha | $dockerSha"

# update only if a saved sha does not equal the latest one
if [ "$savedReSha" != "$reSha" ] || [ "$savedBaseImgSha" != "$baseImgSha" ] || [ "$savedDockerSha" != "$dockerSha" ]; then
    echo "Build $prefix$BRANCH"
    recommit="$reSha" branch="$BRANCH" prefix="$prefix" preimage="$preimage" ./scripts/build-docker.sh
    if [ $? -ne 0 ]; then # dont save sha if something failed
        exit 1
    fi
    # save latest shas
    echo "Save re sha $reSha as $prefix$BRANCH"
    echo "$reSha" > ~/.re-docker/sha/re/"$prefix$BRANCH".sha
    echo "Save base image sha $baseImgSha as $prefix$BRANCH-alpine"
    echo "$baseImgSha" > ~/.re-docker/sha/baseImg/"$prefix$BRANCH"-alpine.sha
    echo "Save dockerfile sha $dockerSha as Dockerfile_$BRANCH"
    echo "$baseImgSha" > ~/.re-docker/sha/dockerfile/"Dockerfile_$BRANCH".sha
else
    echo "Skip $prefix$BRANCH"
fi
exit 0
