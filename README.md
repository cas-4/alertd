# CAS4 alertd

<img src="https://avatars.githubusercontent.com/u/175958109?s=100&v=4" alt="Logo" align="right"/>

This repo refers to a API for alert creation for the
[CAS-4](http://github.com/cas-4) project for the 
[Context Aware System](https://www.unibo.it/en/study/phd-professional-masters-specialisation-schools-and-other-programmes/course-unit-catalogue/course-unit/2023/479036)
class at the [University of Bologna](https://unibo.it).

## Development

You need

- Go `>= v1.23.0`
- Redis `>= 7`

Now you set up some env variables:

- `DEBUG`: if set `=1` you have a more Gin debug logging lines.
- `ADDRESS`: the url of this API.
- `BACKEND_URL`: the url of the backend API.
- `REDIS`: the url of the Redis server.

## Deploy

Fortunately the deployment is automatized by the GitHub Action `cd.yml` which
pushes the latest release version to a [GHCR.io package](https://github.com/cas-4/alertd/pkgs/container/alertd).

A new version is released using

```
./release.sh X.Y
```

Now you just exec

```text
docker pull ghcr.io/cas-4/alertd:latest
```

Or you can build a new image

```text
docker build -t alertd:latest .
docker run Redis
    -e DEBUG=... \
    -e ADDRESS=... \
    -e BACKEND_URL=... \
    -e REDIS=... \
    alertd:latest
```

Or the Docker compose which puts up also the Redis locally.

```text
docker compose up
```
