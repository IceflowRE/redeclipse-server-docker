# Set the base image
FROM debian

LABEL author="Iceflower S"
LABEL maintainer="Iceflower S"
LABEL email="iceflower@gmx.de"
LABEL version="1.0"
LABEL description="Red Eclipse Development Server"

ENV DEBIAN_FRONTEND noninteractive

# Update application repository list
RUN apt-get -qq update && \
    apt-get install --no-install-recommends -y build-essential zlib1g-dev git curl wget ca-certificates cmake pkg-config

# get server source and default maps
RUN git clone -b master https://github.com/red-eclipse/base /temp && \
    git clone -b master https://github.com/red-eclipse/maps.git /redeclipse/data/maps
RUN useradd --create-home --shell /bin/bash redeclipse
RUN chown redeclipse: -R /redeclipse

WORKDIR /temp
USER redeclipse

RUN mkdir build && \
    cd build && \
    cmake ../redeclipse/src -DBUILD_CLIENT=0 -B/redeclipse/bin/amd64 && \
    make clean install -j4

WORKDIR /
COPY /temp/config/ /redeclipse/config/

EXPOSE 28804 28805

CMD /redeclipse/bin/amd64 -h/home/redeclipse/server-config/

