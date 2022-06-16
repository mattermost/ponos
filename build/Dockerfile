# This dockerfile is used to build Mattermost ponos
# A multi stage build, with golang used as a builder
# and gcr.io/distroless/static as runner
ARG GO_IMAGE=golang:1.17@sha256:79138c839452a2a9d767f0bba601bd5f63af4a1d8bb645bf6141bff8f4f33bb8
# hadolint ignore=DL3006
FROM ${GO_IMAGE} as builder

#GO_BUILD_PLATFORMS holds the platforms that we will build the docker image against
ARG GO_BUILD_PLATFORMS=linux-amd64
ARG TERRAFORM_VERSION=1.1.7

# Setup directories structure and compile
COPY . /src
WORKDIR /src

ENV BUILD_SERVICE=server
RUN make go-build
RUN apt-get update -yq && apt-get install --no-install-recommends -yq unzip=6.0-26
# Setup terraform binary
WORKDIR /src/dist
RUN curl -Lo terraform.zip https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip &&\
    unzip -o terraform.zip &&\
    chmod +x terraform

# Shrink final image since we only need the binary
# and use distroless container image as runner for security
FROM gcr.io/distroless/static@sha256:d6fa9db9548b5772860fecddb11d84f9ebd7e0321c0cb3c02870402680cc315f as runner

LABEL name="Ponos" \
    maintainer="cloud-team@mattermost.com" \
    vendor="Mattermost" \
    distribution-scope="public" \
    architecture="x86_64" \
    url="https://mattermost.com" \
    io.k8s.description="Ponos is a toil work elimination tool for SRE tasks." \
    io.k8s.display-name="Mattermost Ponos"

COPY --from=builder /src/dist/server-ponos-linux-amd64 /opt/ponos/bin/ponos-service
COPY --from=builder /src/dist/terraform /usr/local/bin/

# We should refrain from running as privileged user
# Run as UID for nobody
USER 65534

WORKDIR /opt/ponos/bin
ENTRYPOINT ["./ponos-service"]
