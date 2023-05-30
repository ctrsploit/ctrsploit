# syntax=docker/dockerfile:1

ARG GO_VERSION=1.20.4
ARG BASE_DEBIAN_DISTRO="bullseye"
ARG GOLANG_IMAGE="golang:${GO_VERSION}-${BASE_DEBIAN_DISTRO}"
ARG APT_MIRROR

FROM ${GOLANG_IMAGE} AS base
WORKDIR /root/ctrsploit
RUN sed -ri "s/(httpredir|deb).debian.org/${APT_MIRROR:-deb.debian.org}/g" /etc/apt/sources.list \
 && sed -ri "s/(security).debian.org/${APT_MIRROR:-security.debian.org}/g" /etc/apt/sources.list

FROM base AS gox
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
        GOBIN=/build/ GO111MODULE=on go install github.com/mitchellh/gox@latest \
     && /build/gox --help

FROM base AS build
WORKDIR /root/ctrsploit
COPY --from=gox /build/ /usr/local/bin/
RUN --mount=type=cache,sharing=locked,id=moby-build-aptlib,target=/var/lib/apt \
    --mount=type=cache,sharing=locked,id=moby-build-aptcache,target=/var/cache/apt \
        apt update && apt install -y \
            upx
RUN --mount=type=bind,target=.,rw \
    --mount=type=cache,target=/root/.cache/go-build,id=ctrsploit-build \
    --mount=type=cache,target=/go/pkg/mod,id=ctrsploit-mod \
    --mount=type=tmpfs,target=/go/src/ \
    make build-ctrsploit && cp bin/release /build -r
#    GOPROXY=https://goproxy.io,https://goproxy.cn,direct go build -o /build/ -v github.com/ctrsploit/ctrsploit/cmd/ctrsploit

# usage:
# > docker buildx bake binary
# or
# > make binary
FROM scratch AS binary
COPY --from=build /build /

FROM base AS shell
COPY . .