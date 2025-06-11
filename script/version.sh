#!/usr/bin/env bash

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
PROJECT_DIR=$(dirname ${SCRIPT_DIR})
pushd ${PROJECT_DIR} > /dev/null

function get_version() {
    # detect is there any tag
    if git describe --tags --abbrev=0 >/dev/null 2>&1; then
        # get latest tag
        local LATEST_TAG=$(git describe --tags --abbrev=0)
    else
        # setup as the default version if there's no tag
        local LATEST_TAG="init"
    fi

    # get the commits count after the latest tag
    if [ "$LATEST_TAG" == "init" ]; then
        local COMMITS_SINCE_TAG=$(git rev-list HEAD --count)
    else
        local COMMITS_SINCE_TAG=$(git rev-list ${LATEST_TAG}..HEAD --count)
    fi

    if [ "$COMMITS_SINCE_TAG" -eq "0" ]; then
        # there's no new commit after the latest tag
        VERSION=$LATEST_TAG
        RELEASE_DIR="bin/release/${VERSION}"
    else
        # there's new commits after the latest tag
        local COMMIT=$(git rev-parse --short HEAD)
        VERSION="${LATEST_TAG}-${COMMIT}"
        RELEASE_DIR="bin/dev/${VERSION}"
    fi

    # there's changes without store, or stored but not committed
    if ! git diff --quiet || ! git diff --cached --quiet || git ls-files --others --exclude-standard | grep -q .; then
        VERSION="${VERSION}-dirty"
        RELEASE_DIR="bin/dev/${VERSION}"
    fi
}

get_version
echo "$VERSION"
popd > /dev/null
