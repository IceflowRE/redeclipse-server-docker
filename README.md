# Red Eclipse Server Docker
[![Build Status](https://travis-ci.org/IceflowRE/re-server_docker_test.svg?branch=master)](https://travis-ci.org/IceflowRE/re-server_docker_test)
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
---

This provides a Dockerfile for an easy handling and maintaining of a [Red Eclipse](https://redeclipse.net/) Server.
Currently the Docker images are build against the latest commits and will be checked for updates once a day.

### Images
Latest images are available at [Docker Hub](https://hub.docker.com/r/iceflower/red-eclipse_devel_server_test/).
Pull them with `docker pull iceflower/red-eclipse_devel_server_test:<tag>`.

#### Tags
|       Tag      | Server type | Architecture |
|:--------------:|:-----------:|:------------:|
|     master     | development |     amd64    |
| arm64v8-master | development |    arm64v8   |
|     stable     |    stable   |     amd64    |
| arm64v8-stable |    stable   |    arm64v8   |

**Note: arm64v8 images are NOT available at the moment and will follow later!**

### Usage
- Pull latest docker image from Docker Hub.
`docker pull iceflower/red-eclipse_devel_server_test:<tag>`

  - `<tag>`: an available image tag
- Run them.  
`docker run -p <serverport>:<serverport>/udp -p <serverport + 1>:<serverport + 1>/udp -v <re home dir>:/re-server-config/home -v <re package dir>:/re-server-config/package -v <sauerbraten maps dir>:/re-server-config/sauer -v <log dir>:/home/redeclipse/re-log iceflower/red-eclipse_devel_server_test:<tag>`

  - `<serverport>:` the serverport specified inside the `servinit.cfg` from your server
  - `<serverport + 1>`: the serverport + 1  
  ***you can link host directories inside the docker container, if dont want to link a directory just leave off the specific `-v` parameter.***
  - `<re home dir>`: RE home directory on your host system, **must linked**
  - `<re package dir>`: package directory, mostly maps, on your host system
  - `<sauerbraten maps dir>`: sauerbraten map directory on your host system
  - `<log dir>`: log directory on your host system
  - `<tag>`: an available image tag

#### Example
`docker run -p 28803:28803/udp -p 28804:28804/udp -v /home/iceflower/redeclipse-config/devel_home:/re-server-config/home -v /home/iceflower/redeclipse-config/package:/re-server-config/package -v /home/iceflower/redeclipse-config/sauerbraten:/re-server-config/sauer -v /home/iceflower/redeclipse-config/logs/devel_log:/home/redeclipse/re-log iceflower/red-eclipse_devel_server_test`

---

## Web
https://github.com/IceflowRE/re-server_docker_test

## Credits
- Developer
  - Iceflower S

## License
![Image of GPLv3](http://www.gnu.org/graphics/gplv3-127x51.png)

Copyright 2017 - now ©  Iceflower S

This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.
This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
You should have received a copy of the GNU General Public License along with this program; if not, see <http://www.gnu.org/licenses/gpl.html>.
