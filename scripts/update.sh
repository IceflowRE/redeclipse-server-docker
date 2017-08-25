#! /bin/bash

# variable BRANCH is defined inside the travis.yml env.
SHA="$(git ls-remote --heads https://github.com/red-eclipse/base.git master | awk '{ print $1 }')"

if [ "$(cat ./sha/$BRANCH.sha)" != "$SHA" ]; then
    echo "Build $BRANCH"
    ./scripts/build-docker.sh "$BRANCH"
# dont save sha if something failed
    if [ $? -ne 0 ]; then
        exit 1
    fi
    echo "Save sha $SHA as $BRANCH"
    echo "$SHA" > "./sha/$BRANCH.sha"
else
    echo "Skip $BRANCH"
fi
exit 0
