#! /bin/bash
# execute from project root

# master
docker manifest create iceflower/redeclipse-server:master iceflower/redeclipse-server:amd64-master iceflower/redeclipse-server:arm64v8-master
docker manifest annotate iceflower/redeclipse-server:master iceflower/redeclipse-server:arm64v8-master --variant armv8
docker manifest push --purge iceflower/redeclipse-server:master
if [ $? -ne 0 ]; then
    exit 1
fi

# stable
docker manifest create iceflower/redeclipse-server:stable iceflower/redeclipse-server:amd64-stable iceflower/redeclipse-server:arm64v8-stable
docker manifest annotate iceflower/redeclipse-server:stable iceflower/redeclipse-server:arm64v8-stable --variant armv8
docker manifest push --purge iceflower/redeclipse-server:stable
if [ $? -ne 0 ]; then
    exit 1
fi

exit 0
