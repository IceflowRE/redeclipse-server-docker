#! /bin/bash

# variable BRANCH is defined inside the travis.yml env.

if [ $TRAVIS_TEST_RESULT ]; then
    echo "Save sha $SHA as $BRANCH"
    echo "$SHA" > "./sha/$BRANCH.sha"
fi
