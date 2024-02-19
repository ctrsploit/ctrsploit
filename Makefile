.PHONY: all shell local build

APP_NAME := ctrsploit

# mirror
DEFAULT_CN_APT_MIRROR := mirrors.tuna.tsinghua.edu.cn
DEFAULT_CN_GOPROXY := https://goproxy.cn,https://goproxy.io,direct
APT_MIRROR ?= $(if $(CN),$(DEFAULT_CN_APT_MIRROR),)
GOPROXY ?= $(if $(CN),$(DEFAULT_CN_GOPROXY),)

# debug
PROGRESS_PLAIN := --progress plain
DEBUG_FLAGS ?= $(if $(DEBUG),$(PROGRESS_PLAIN),)

# ldflags
GIT_COMMIT := $(shell git rev-parse --short HEAD || echo unsupported)
VERSION := $(shell cat ./VERSION)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
SLIM_LDFLAGS ?= -s -w
LDFLAGS := "$(SLIM_LDFLAGS) \
	-X github.com/ctrsploit/sploit-spec/pkg/version.Version=${VERSION} \
	-X github.com/ctrsploit/sploit-spec/pkg/version.GitCommit=${GIT_COMMIT} \
	-X github.com/ctrsploit/sploit-spec/pkg/version.BuildTime=${BUILD_TIME}"

# image
DEV_IMAGE := ${APP_NAME}-dev
DOCKERFILE := Dockerfile_dev

# build flags
BUILD_APT_MIRROR := $(if $(APT_MIRROR),--build-arg APT_MIRROR="$(APT_MIRROR)")
BUILD_GO_PROXY := $(if $(GOPROXY),--build-arg GOPROXY="$(GOPROXY)")
BUILD_SLIM_LDFLAGS := $(if $(SLIM_LDFLAGS), --build-arg SLIM_LDFLAGS="$(SLIM_LDFLAGS)")
BUILD_OPTS := ${BUILD_APT_MIRROR} ${BUILD_GO_PROXY} ${BUILD_SLIM_LDFLAGS} ${DOCKER_BUILD_ARGS} ${DOCKER_BUILD_OPTS} -f "$(DOCKERFILE)"

binary: bundle
	APT_MIRROR="$(APT_MIRROR)" GOPROXY="$(GOPROXY)" SLIM_LDFLAGS="$(SLIM_LDFLAGS)" docker buildx bake binary ${DEBUG_FLGAS}

bundle:
	mkdir -p bin/release

build:
	LDFLAGS=$(LDFLAGS) ./release.sh

install: build
	rm -f /usr/local/bin/${APP_NAME} && ln -s $(CURDIR)/bin/release/${APP_NAME}_linux_amd64 /usr/local/bin/${APP_NAME}

image:
	docker buildx build $(BUILD_OPTS) --load -t "$(DEV_IMAGE)" ${DEBUG_FLAGS} .

shell: image
	docker run --rm -ti -v $(CURDIR):/root/app $(DEV_IMAGE) bash

# usage:
# make binary
# make shell
# make install

# args:
# SLIM_FLAGS=
#	default to -s -w, setup to empty to disable slim ldflags
# CN=1
#	use cn mirrors
# DEBUG=1
#	build --progress=plain
