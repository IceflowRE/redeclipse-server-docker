#! /bin/bash

# variable BRANCH is defined inside the travis.yml env.
SHA="$(git ls-remote --heads https://github.com/red-eclipse/base.git master | awk '{ print $1 }')"

if [ "$(cat ./sha/$BRANCH.sha)" != "$SHA" ]; then
    echo "Build $BRANCH"
    ./scripts/build-docker.sh "$BRANCH"
else
    echo "Skip $BRANCH"
fi
