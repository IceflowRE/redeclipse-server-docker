# Red Eclipse Server Docker

[![maintained](https://img.shields.io/badge/maintained-yes-brightgreen.svg)][github]
[![][github actions images]][github actions]
[![DockerHub](https://img.shields.io/badge/Docker_Hub--FF69A4.svg?style=social)][docker hub]
[![Github](https://img.shields.io/badge/Github--FF69A4.svg?style=social)][github]

---

This provides the source for an easy handling and maintaining Docker image of a [Red Eclipse](https://redeclipse.net/)
Server.  
Additional with a Go console application which will update the DockerHub images.  
Currently the Docker images are build against the latest commits and will be checked for updates once a day.

---

## Supported tags and respective `Dockerfile` links

- [`master`](https://github.com/IceflowRE/redeclipse-server-docker/blob/main/Dockerfile_master)
- [`stable`](https://github.com/IceflowRE/redeclipse-server-docker/blob/main/Dockerfile_stable)
- [`v1.5.3`](https://github.com/IceflowRE/redeclipse-server-docker/blob/main/Dockerfile_stable)
- [`v1.5.5`](https://github.com/IceflowRE/redeclipse-server-docker/blob/main/Dockerfile_stable)
- [`v1.5.6`](https://github.com/IceflowRE/redeclipse-server-docker/blob/main/Dockerfile_stable)
- [`v1.5.8`](https://github.com/IceflowRE/redeclipse-server-docker/blob/main/Dockerfile_stable)
- [`v1.6.0`](https://github.com/IceflowRE/redeclipse-server-docker/blob/main/Dockerfile_stable)
- [`v2.0.0`](https://github.com/IceflowRE/redeclipse-server-docker/blob/main/Dockerfile_2_0_0)

\* *`stable` does not mark the latest stable release, it tags the latest legacy version (v1.x.x)*

Available architectures are `amd64`.

If your architecture is not available check if RE can be build on that architecture and follow [An image is not available for my architecture?!
](#An image is not available for my architecture?!).

## How to use this image

We will use docker compose as it is the easiest way to manage the running container.

Create a Docker Compose file and name it `docker-compose.yml`, as a starting point you can use [docker-compose-template.yml](https://github.com/IceflowRE/redeclipse-server-docker/blob/main/docker-compose-template.yml)

Replace the `<variable>` including the brackets.

```yml
services:
    # <service_name>: the name to access this later (e.g. `master`, `v2_0_0`)
    <service_name>:
        # <tag>: the image tag you want to use (e.g. `master`, `v2.0.0`)
        image: iceflower/redeclipse-server:<tag>
        ports:
            # <serverport>: this port will be published and accessible from outside,
            # the port number must match port defined in RE's `servinit.cfg`
            - "<serverport>:<serverport>/udp"
            # <serverport + 1>: the server port above plus one
            - "<serverport + 1>:<serverport + 1>/udp"
        restart: unless-stopped
        volumes:
            # <RE home dir>: path to the RE home/ config directory on your host system 
            # (e.g. `/home/iceflower/re-master/home`)
            - type: bind
              source: <RE home dir>
              target: /re-server-config/home
              read_only: true
            # <RE package dir>: path to the RE package directory on your host system, you can place custom maps there
            # if you do not want this, just remove this section (e.g. `/home/iceflower/re-master/package`)
            - type: bind
              source: /home/iceflower/redeclipse-config/package
              target: /re-server-config/package
              read_only: true
            # <sauerbraten dir>: path to a Sauerbraten directory/installation
            # if you use a version higher or equal `v2.0.0` or `master` remove this section
            # (e.g. `/home/iceflower/sauerbraten`)
            - type: bind
              source: <sauerbraten dir>
              target: /re-server-config/sauer
              read_only: true
        logging:
            options:
                max-size: "2000k"
                max-file: "10"
```

To pull/start/stop a specific service add the service name to the end otherwise, it is applied to all.

- Pull the latest docker image from Docker Hub for all defined services  
  `docker compose pull`

- Start/ Restart container  
  `docker compose up -d`

- Shutdown and wait a maximum of 10 minutes before forcing it  
  `docker compose stop --time=600`

#### Multiple Server

Copy and paste the whole service section above and change the values (service name, port, home directory, etc.)

### An image is not available for my architecture?!

Follow the table below to copy the required `Dockerfile` next to your `docker-compose.yml`. Have the chosen `Dockerfile` in mind.

| Version  |     Dockerfile      |
|:--------:|:-------------------:|
|  master  | `Dockerfile_master` |
|  stable  | `Dockerfile_stable` |
| < v2.0.0 | `Dockerfile_stable` |
|  v2.0.0  | `Dockerfile_2_0_0`  |

Create a file named `.dockerignore` next to the `Dockerfile` with the content `**`.

Edit your `docker-compose-yml` and replace the `image: ...` part in the service you want to build with

```yml
build:
    # <dockerfile>: name of the dockerfile, usable for the chosen tag below
    dockerfile: <dockerfile>
    args:
        # <tag>: can be any git reference (branch name, tag, SHA) (e.g. `master`, `v2.0.0`, ...)
        TAG: <tag>
        RE_COMMIT: ""
        ALPINE_SHA: ""
        DOCKERFILE_SHA: ""
```

Then use `docker compose build` to build the custom image.

## Update server automatically

Create a file with the name `update-server.sh` (make sure it has executable rights).
You can place it right next to the `docker-compose.yml`.

Adjust the directory path according where you place the file.

```bash
#!/bin/bash -e

# EDIT NEXT LINE
cd /home/iceflower/re/
docker compose stop -t 600
# docker compose build
docker compose pull
docker compose up -d
```

To update regularly you can create a cron job, how to do this exactly refer to a guide matching your OS.

```cron
0 3 * * * /home/iceflower/re/update-server.sh > /home/iceflower/re/cron.log 2>&1
```

This will update every day at 3:00.

## DockerHub Image Updater

Build with

```shell
cd go-docker-updater
go build -x -o updater ./cmd/updater/
```

For more options see `--help`.

The updater is used to update the DockerHub images.

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
