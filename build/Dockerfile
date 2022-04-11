# Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
# See LICENSE.txt for license information.

# Build the mattermost ponos app
ARG DOCKER_BUILD_IMAGE=golang:1.16
ARG DOCKER_BASE_IMAGE=alpine:3.14
ARG BUILD_VERSION

FROM ${DOCKER_BUILD_IMAGE} AS build
WORKDIR /opt/ponos
COPY . /opt/ponos
RUN make build-linux

# Final Image
FROM ${DOCKER_BASE_IMAGE}

LABEL name="Ponos" \
    maintainer="cloud-team@mattermost.com" \
    vendor="Mattermost" \
    distribution-scope="public" \
    architecture="x86_64" \
    url="https://mattermost.com" \
    io.k8s.description="Ponos is a toil work elimination tool for SRE tasks." \
    io.k8s.display-name="Mattermost Ponos"

ENV PONOS=/ponos/mattermost-app-ponos \
    USER_UID=10001 \
    USER_NAME=ponos \
    PONOS_DIR=/ponos

RUN apk update && apk add libc6-compat && apk add ca-certificates
COPY --from=build /opt/ponos/build/_output/bin/mattermost-app-ponos-linux-amd64 /chaos/mattermost-app-ponos
COPY --from=build /opt/ponos/build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

WORKDIR /ponos

USER ${USER_UID}

EXPOSE  3000

ENTRYPOINT ["/usr/local/bin/entrypoint"]
