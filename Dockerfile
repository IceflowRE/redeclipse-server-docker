# Set the base image
FROM alpine

# informations
LABEL author="Iceflower S"
LABEL maintainer="Iceflower S"
LABEL email="iceflower@gmx.de"
LABEL version="1.0"
LABEL description="Red Eclipse Development Server"

# Add server user and set permissions
RUN adduser -S -D redeclipse \
    && mkdir /redeclipse \
    && chown redeclipse: -R /redeclipse

# Update application repository list, create build dir, build server, move server files, create other permanent files and clean up
RUN apk update \
    && apk add --no-cache --virtual build-deps gcc g++ zlib-dev git ca-certificates coreutils cmake make \
    && apk add --no-cache libstdc++ \
    && git clone -b master https://github.com/red-eclipse/base /temp \
    \
    && mkdir /temp/build \
    && cd /temp/build \
    && cmake ../src -DBUILD_CLIENT=0 \
    && make clean install -j4 \
    \
    && mkdir -p /redeclipse/config \
    && mv /temp/config/ /redeclipse/ \
    && mkdir -p /redeclipse/bin/amd64 \
    && mv /temp/bin/amd64/redeclipse_server_linux /redeclipse/bin/amd64/redeclipse_server_linux \
    \
    && apk update \
    && apk del build-deps \
    && rm -rf /temp

# Add defaults maps and server config folder
RUN apk add --no-cache --virtual deps git ca-certificates \
    && git clone -b master https://github.com/red-eclipse/maps.git /redeclipse/data/maps \
    && mkdir -p /re-server-config/home \
    && mkdir -p /re-server-config/package \
    && mkdir -p /re-server-config/sauer \
    && apk update \
    && apk del deps

USER redeclipse

# This ports have to be used by the server config
EXPOSE 28800/udp 28801/udp

CMD cd /redeclipse && ./bin/amd64/redeclipse_server_linux -h/re-server-config/home -p/re-server-config/package -o/re-server-config/sauer
