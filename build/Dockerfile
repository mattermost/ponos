# Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
# See LICENSE.txt for license information.

# Build the mattermost ponos server
ARG DOCKER_BUILD_IMAGE=golang:1.17.9@sha256:55759506d9a7b33b28117977170db9713903799c110db99af2c15a3f603602af
ARG DOCKER_BASE_IMAGE=gcr.io/distroless/static@sha256:d6fa9db9548b5772860fecddb11d84f9ebd7e0321c0cb3c02870402680cc315f
ARG BUILD_VERSION

FROM ${DOCKER_BUILD_IMAGE} AS build
WORKDIR /opt/ponos
COPY . /opt/ponos
RUN apt-get update -yq && apt-get install -yq unzip
RUN make get-terraform

ENV BUILD_SERVICE=server
RUN make build

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

COPY --from=build /opt/ponos/build/_output/bin/ponos-server-linux-amd64 /ponos/ponos-service
COPY --from=build /opt/ponos/build/terraform /usr/local/bin/
COPY --from=build /opt/ponos/terraform /ponos/terraform

# Run as UID for nobody
USER 65534

WORKDIR /ponos
ENTRYPOINT ["./ponos-service"]
