#! /bin/bash

# variable BRANCH is defined inside the travis.yml env.
sha="$(git ls-remote --heads https://github.com/red-eclipse/base.git $BRANCH | awk '{ print $1 }')"

if [ "$(cat ./sha/$BRANCH.sha)" != "$sha" ]; then
    echo "Build $BRANCH"
    ./scripts/build-docker.sh "$BRANCH"
# dont save sha if something failed
    if [ $? -ne 0 ]; then
        exit 1
    fi
    echo "Save sha $sha as $BRANCH"
    echo "$sha" > "./sha/$BRANCH.sha"
else
    echo "Skip $BRANCH - $sha"
fi
exit 0
