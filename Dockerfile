# Set the base image
FROM debian
# Dockerfile author / maintainer 
MAINTAINER Iceflower S <iceflower@gmx.de>

ENV DEBIAN_FRONTEND noninteractive

# Update application repository list
RUN apt-get -qq update && \
    apt-get install --no-install-recommends -y build-essential zlib1g-dev git curl wget ca-certificates cmake pkg-config

# get server source and default maps
RUN git clone -b master https://github.com/red-eclipse/base redeclipse && \
    cd redeclipse/data && \
    git submodule update --init maps && \
    cd ../..
RUN useradd --create-home --shell /bin/bash redeclipse

WORKDIR /redeclipse
RUN chown redeclipse: -R /redeclipse

USER redeclipse
RUN mkdir build && \
    cd build && \
    cmake ../src -DBUILD_CLIENT=0 && \
    make clean install -j4

EXPOSE 28804 28805

CMD /redeclipse/redeclipse_server.sh

