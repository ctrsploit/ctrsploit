# syntax=docker/dockerfile:1

ARG GO_VERSION=1.20.4
ARG BASE_DEBIAN_DISTRO="bullseye"
ARG GOLANG_IMAGE="golang:${GO_VERSION}-${BASE_DEBIAN_DISTRO}"

FROM ${GOLANG_IMAGE} AS base
ARG APT_MIRROR
WORKDIR /root/app
RUN sed -ri "s/(httpredir|deb).debian.org/${APT_MIRROR:-deb.debian.org}/g" /etc/apt/sources.list \
 && sed -ri "s/(security).debian.org/${APT_MIRROR:-security.debian.org}/g" /etc/apt/sources.list \
 && sed -ri "s/(snapshot).debian.org/${APT_MIRROR:-snapshot.debian.org}/g" /etc/apt/sources.list \
 && cat /etc/apt/sources.list
RUN git config --global --add safe.directory /root/app

FROM base AS gox
ARG GOPROXY
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
        GOBIN=/build/ GO111MODULE=on GOPROXY=${GOPROXY} go install github.com/mitchellh/gox@latest \
     && /build/gox --help

FROM base AS build-env
ARG GOPROXY
WORKDIR /root/app
COPY --from=gox /build/ /usr/local/bin/
RUN --mount=type=cache,sharing=locked,id=ctrsploit-build-aptlib,target=/var/lib/apt \
    --mount=type=cache,sharing=locked,id=ctrsploit-build-aptcache,target=/var/cache/apt \
        apt update && apt install -y \
            upx jq
RUN --mount=type=bind,target=.,rw \
    --mount=type=cache,target=/root/.cache/go-build,id=ctrsploit-build \
    --mount=type=cache,target=/go/pkg/mod,id=ctrsploit-mod \
    --mount=type=tmpfs,target=/go/src/ \
    GOPROXY=${GOPROXY} go mod download

FROM build-env AS build
RUN --mount=type=bind,target=.,rw \
    --mount=type=cache,target=/root/.cache/go-build,id=ctrsploit-build \
    --mount=type=cache,target=/go/pkg/mod,id=ctrsploit-mod \
    --mount=type=tmpfs,target=/go/src/ \
    make build-ctrsploit && mv bin/release /build

# usage:
# > docker buildx bake binary
# or
# > make binary
FROM scratch AS binary
COPY --from=build /build /

FROM build-env AS shell
ARG GOPROXY
ENV GOPROXY=${GOPROXY}
