# Ponos

## Background

[πόνος](https://en.wikipedia.org/wiki/Ponos) - pain in greek!

## Introduction

Ponos is tool which the Mattermost SRE team uses daily to eliminate toil work with ChatOps tools. It includes the followings:
- Ponos microservice: it's responsible for the business logic of the toil work
- Ponos Mattemost App: it's the ChatOps based on [Mattermost App Framework](https://developers.mattermost.com/integrate/apps/) which interacts with Ponos service.

## Developing

### Running

For Ponos HTTP service you need a running provisioner and management of workspaces:

```bash
export PONOS_PROVISIONER_ADDRESS=<your-provisioner-url>
export PONOS_WORKSPACES_ADDRESS=<your-workspaces-url>
export BUILD_SERVICE=server
make run
```

For Ponos `ChatOps` Mattermost App you need to run the following:

```bash
export PONOS_PROVISIONER_ADDRESS=<your-provisioner-url>
export PONOS_WORKSPACES_ADDRESS=<your-workspaces-url>
export BUILD_SERVICE=app
make run
```

### Testing

Running tests:

```bash
make test
```

### Docker

Docker build and run locally Ponos Service.

```bash
make docker-build

docker run --rm -p 3000:3000 \
    -e PONOS_PROVISIONER_ADDRESS=https://<your-url> \
    -e PONOS_WORKSPACES_ADDRESS=https://<your-url>  mattermost/ponos-service:dev-local
```
