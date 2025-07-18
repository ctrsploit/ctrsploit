#!/usr/bin/env bash

PKG=$1
if [[ -z "$PKG" ]]; then
  echo "Usage: $0 <package>"
  exit 1
fi

SCRIPT_DIR=$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")
PROJECT_DIR=$(dirname ${SCRIPT_DIR})
pushd ${PROJECT_DIR} > /dev/null

TEST_BINARY_NAME="${PKG//\//_}.test"
mkdir -p bin/test
rm -f bin/test/${TEST_BINARY_NAME}
CGO_ENABLED=0 go test -c -o bin/test/${TEST_BINARY_NAME} ${PKG}

popd > /dev/null