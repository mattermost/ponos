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
export AWS_REGION=us-east-1
export AWS_PROFILE="<profile-for-infra-migrations>"
export PONOS_ACCOUNT_ID="<aws-account-id-for-testing>"
export PONOS_KMS_KEY="<aws-kms-id-for-testing>"
export PONOS_PROVISIONER_ADDRESS=<your-provisioner-url>
export PONOS_WORKSPACES_ADDRESS=<your-workspaces-url>
export BUILD_SERVICE=server
export PONOS_DB_DSN="host=<db host> user=<db user> password=<db user's password> dbname=<db name>"
make run
```

For local development, you can run a local PostgreSQL database via

```bash
docker-compose up
```

by using the following as the databases's DSN:

```
export PONOS_DB_DSN="host=localhost user=ponos_dev_user password=ponos_dev_password dbname=ponos_dev"
```

For Ponos `ChatOps` Mattermost App you need to run the following:

```bash
export PONOS_APP_ROOT_URL=http://<your-ip>:3000
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
