# Red Eclipse Server Docker
[![maintained](https://img.shields.io/badge/maintained-yes-brightgreen.svg)][github]
[![DockerHub](https://img.shields.io/badge/Docker_Hub--FF69A4.svg?style=social)][docker hub]
[![stable](https://badges.herokuapp.com/travis.com/IceflowRE/redeclipse-server-docker?env=BRANCH=stable&label=stable)][travis ci]
[![master](https://badges.herokuapp.com/travis.com/IceflowRE/redeclipse-server-docker?env=BRANCH=master&label=master)][travis ci]
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
---

This provides the source for an easy handling and maintaining Docker image of a [Red Eclipse](https://redeclipse.net/) Server.  
Additional with an python console application which can be run automatically to update it.  
Currently the Docker images are build against the latest commits and will be checked for updates once a day.

---

## Images
Latest images are available at [Docker Hub][docker hub].  
Pull them with `docker pull iceflower/redeclipse-server:<tag>`. Use as tag `master` or `stable`, the correct architexture will be choosen.  
The `amd64` images are build with [Travis CI][travis ci]. The `arm64v8` images are build on an own server.

|       Tag      | Server type | Architecture |              Build              |                  Size / Layers                |
|:--------------:|:-----------:|:------------:|:-------------------------------:|:---------------------------------------------:|
|  amd64-stable  |    stable   |     amd64    | [![][travis stable]][travis ci] |     [![][mbadge stable]][mbadge stable l]     |
| arm64v8-stable |    stable   |    arm64v8   |     [![][no build]][github]     | [![][mbadge arm stable]][mbadge arm stable l] |
|  amd64-master  | development |     amd64    | [![][travis master]][travis ci] |     [![][mbadge master]][mbadge master l]     |
| arm64v8-master | development |    arm64v8   |     [![][no build]][github]     | [![][mbadge arm master]][mbadge arm master l] |

## Usage
Replace the variables with the respective values.

  - `<name>` a container name, under which it will run
  - `<tag>` an available image tag (either `master` or `stable`)
  - `<serverport>` the serverport specified inside the `servinit.cfg` from your server
  - `<serverport + 1>` the serverport + 1  
  ***you can link host directories inside the docker container, if dont want to link a directory just leave off the specific `-v` parameter.***
  - `<re home dir>` RE home directory on your host system, **must linked**
  - `<re package dir>` package directory, inside you can place a maps directory, on your host system
  - `<sauerbraten dir>` sauerbraten directory on your host system

### Docker Compose (recommend)
- Create own Docker Compose file, as base you can use [docker-compose.yml.template](./docker-compose.yml.template)  
  - *Create a copy with name `docker-compose.yml`*
  - *Change all the `<variable>` inside the file, to their respective values*

- Pull latest docker image from Docker Hub for all defined services  
`docker-compose pull`

- Start/ Restart container (for all specified services, dont write any name)  
`docker-compose -p re-server up -d <name>`

- Shutdown and wait a maximum of 10 seconds before forcing (for all specified services, dont write any name)  
`docker-compose stop <name>`

#### Multiple Server
Copy and paste the whole section below the point `services` and change the values. Then start it with the new other name.

#### Example
- Create own Docker Compose file, as base you can use [docker-compose.yml.template](./docker-compose.yml.template)  
[docker-compose.yml.example](./docker-compose.yml.example)

- Pull latest docker image from Docker Hub for all defined services  
`docker-compose pull`

- Start/ Restart container  
`docker-compose -p re-server up -d master`

- Shutdown and wait a maximum of 10 minutes before forcing  
`docker-compose stop --time=600 master`

### Command line
- Pull latest docker image from Docker Hub  
`docker pull iceflower/redeclipse-server:<tag>`

- Create container  
```
docker create \
-p <serverport>:<serverport>/udp \
-p <serverport + 1>:<serverport + 1>/udp \
--mount type=bind,source="<re home dir>",target=/re-server-config/home,readonly=true \
--mount type=bind,source="<re package dir>",target=/re-server-config/package,readonly=true \
--mount type=bind,source="<sauerbraten dir>",target=/re-server-config/sauer,readonly=true \
--name <name> \
iceflower/redeclipse-server:<tag>
```

- Start container  
`docker start <name>`

- Shutdown and wait a maximum of 10 seconds before forcing  
`docker stop <name>`

#### Multiple Server
Create a container, with changed values and another name and start it.

#### Example
- Pull latest docker image from Docker Hub  
`docker pull iceflower/redeclipse-server:master`

- Create container  
```
docker create \
-p <serverport>:<serverport>/udp \
-p <serverport + 1>:<serverport + 1>/udp \
--mount type=bind,source="/home/iceflower/redeclipse-config/devel_home",target=/re-server-config/home,readonly=true \
--mount type=bind,source="/home/iceflower/redeclipse-config/package",target=/re-server-config/package,readonly=true \
--mount type=bind,source="/home/iceflower/redeclipse-config/sauerbraten",target=/re-server-config/sauer,readonly=true \
--name re-server-master \
iceflower/redeclipse-server:master
```

- Start container  
`docker start re-server-master`

- Shutdown and wait a maximum of 10 minutes before forcing  
`docker stop --time=600 re-server-master`

---

## Web
https://github.com/IceflowRE/redeclipse-server-docker

## Credits
- Developer
  - Iceflower S
    - iceflower@gmx.de

### External Tools
- badge-matrix *by* Brian Beck ([exogen](https://github.com/exogen))
    - https://github.com/exogen/badge-matrix
    - MIT License

## License
![Image of GPLv3](http://www.gnu.org/graphics/gplv3-127x51.png)

Copyright 2017 - now © Iceflower S

This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program; if not, see <http://www.gnu.org/licenses/gpl.html>.

---

**The MIT License** *(only for the docker-compose files)*

Copyright 2017 - now Iceflower S

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

[travis ci]: https://travis-ci.com/IceflowRE/redeclipse-server-docker
[github]: https://github.com/IceflowRE/redeclipse-server-docker
[docker hub]: https://hub.docker.com/r/iceflower/redeclipse-server
[no build]: https://img.shields.io/badge/build-inaccessible-lightgrey.svg
[travis stable]: https://badges.herokuapp.com/travis/IceflowRE/redeclipse-server-docker?env=BRANCH=stable&label=build
[travis master]: https://badges.herokuapp.com/travis/IceflowRE/redeclipse-server-docker?env=BRANCH=master&label=build
[mbadge stable]: https://images.microbadger.com/badges/image/iceflower/redeclipse-server:amd64-stable.svg
[mbadge stable l]: https://microbadger.com/images/iceflower/redeclipse-server:amd64-stable
[mbadge master]: https://images.microbadger.com/badges/image/iceflower/redeclipse-server:amd64-master.svg
[mbadge master l]: https://microbadger.com/images/iceflower/redeclipse-server:amd64-master
[mbadge arm stable]: https://images.microbadger.com/badges/image/iceflower/redeclipse-server:arm64v8-stable.svg
[mbadge arm stable l]: https://microbadger.com/images/iceflower/redeclipse-server:arm64v8-stable
[mbadge arm master]: https://images.microbadger.com/badges/image/iceflower/redeclipse-server:arm64v8-master.svg
[mbadge arm master l]: https://microbadger.com/images/iceflower/redeclipse-server:arm64v8-master
