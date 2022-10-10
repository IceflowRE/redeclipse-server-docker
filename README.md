# Red Eclipse Server Docker

[![maintained](https://img.shields.io/badge/maintained-yes-brightgreen.svg)][github]
[![][github actions images]][github actions]
[![DockerHub](https://img.shields.io/badge/Docker_Hub--FF69A4.svg?style=social)][docker hub]

---

This provides the source for an easy handling and maintaining Docker image of a [Red Eclipse](https://redeclipse.net/)
Server.  
Additional with an Go console application which can be run automatically to update the images.  
Currently the Docker images are build against the latest commits and will be checked for updates once a day.

---

## Images

Latest images are available at [Docker Hub][docker hub].  
Pull them with `docker pull iceflower/redeclipse-server:<tag>`.
All images are available for `amd64` and `arm64/v8`, `amd64` images are build with [GitHub Actions][github actions]
, `arm64/v8` images are build on an own server.

| Arch  |                    Status                    |
|:-----:|:--------------------------------------------:|
| amd64 | [![][github actions images]][github actions] |
| arm64 |           [![][no build]][github]            |

Available tags are

- `master` - master branch
- `stable` - stable branch
- `v1.5.3`
- `v1.5.5`
- `v1.5.6`
- `v1.5.8`
- `v1.6.0`
- `v2.0.0`

## Usage

Replace the variables with the respective values.

- `<name>` a container name, under which it will run
- `<tag>` an available image tag (either `master`, `stable-re2` or `stable`)
- `<serverport>` the serverport specified inside the `servinit.cfg` from your server
- `<serverport + 1>` the serverport + 1  
  ***you can link host directories inside the docker container, if dont want to link a directory just leave off the
  specific `-v` parameter.***
- `<re home dir>` RE home directory on your host system, **must linked**
- `<re package dir>` package directory, inside you can place a maps directory, on your host system
- `<sauerbraten dir>` sauerbraten directory on your host system (only available for < v2.0.0)

### Build image yourself

Edit the build section in `docker-compose.yml` and use the correct `dockerfile` and `TAG` (replace the X.X.X with the
version you want).

| Version  |     Dockerfile      |  Tag   |
|:--------:|:-------------------:|:------:|
|  master  | `Dockerfile_master` | master |
|  stable  | `Dockerfile_stable` | stable |
|  v2.0.0  | `Dockerfile_2_0_0`  | v2.0.0 |
| < v2.0.0 | `Dockerfile_stable` | vX.X.X |

Then use `docker compose build`

### Docker Compose (recommend)

- Create own Docker Compose file, as base you can use [docker-compose.yml.template](./docker-compose.yml.template)
    - *Create a copy with name `docker-compose.yml`*
    - *Change all the `<variable>` inside the file, to their respective values*

- Pull latest docker image from Docker Hub for all defined services  
  `docker compose pull`

- Start/ Restart container (for all specified services, dont write any name)  
  `docker compose -p re-server up -d <name>`

- Shutdown and wait a maximum of 10 seconds before forcing (for all specified services, dont write any name)  
  `docker compose stop <name>`

#### Multiple Server

Copy and paste the whole section below the point `services` and change the values. Then start it with the new other
name.

#### Example

- Create own Docker Compose file, as base you can use [docker-compose.yml.template](./docker-compose.yml.template)  
  [docker-compose.yml.example](./docker-compose.yml.example)

- Pull latest docker image from Docker Hub for all defined services  
  `docker compose pull`

- Start/ Restart container  
  `docker compose -p re-server up -d master`

- Shutdown and wait a maximum of 10 minutes before forcing  
  `docker compose stop --time=600 master`

## Updater

Build with

```shell
cd go-docker-updater
go build -x -o updater ./cmd/updater/
```

For more options see `--help`.

The updater can be used to update the docker images.

## Web

https://github.com/IceflowRE/redeclipse-server-docker

## Credits

- Developer
    - Iceflower S
        - iceflower@gmx.de

## License

Copyright 2017-preset Iceflower S

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

---

The Red Eclipse server files are
using [THE RED ECLIPSE LICENSE](https://github.com/redeclipse/base/blob/master/doc/license.txt).

[github actions]: https://github.com/IceflowRE/redeclipse-server-docker/actions/workflows/update_docker_images.yml

[github actions images]: https://img.shields.io/github/workflow/status/IceflowRE/redeclipse-server-docker/Update%20Docker%20images

[github]: https://github.com/IceflowRE/redeclipse-server-docker

[docker hub]: https://hub.docker.com/r/iceflower/redeclipse-server

[no build]: https://img.shields.io/badge/build-inaccessible-lightgrey.svg
