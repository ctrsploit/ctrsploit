#!/bin/bash
set -ex

SCRIPT_DIR=$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")
PROJECT_DIR=$(dirname ${SCRIPT_DIR})

# setup RELEASE_DIR
source ${SCRIPT_DIR}/version.sh
cd ${PROJECT_DIR}
mkdir -p ${RELEASE_DIR}
rm -f bin/latest & ln -s $(realpath --relative-to=bin ${RELEASE_DIR}) bin/latest
cd ${RELEASE_DIR}

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ctrsploit_linux_amd64 -ldflags "${LDFLAGS}" github.com/ctrsploit/ctrsploit/cmd/ctrsploit
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o checksec_linux_amd64 -ldflags "${LDFLAGS}" github.com/ctrsploit/ctrsploit/cmd/checksec
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o env_linux_amd64 -ldflags "${LDFLAGS}" github.com/ctrsploit/ctrsploit/cmd/env

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ctrsploit_linux_arm64 -ldflags "${LDFLAGS}" github.com/ctrsploit/ctrsploit/cmd/ctrsploit
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o checksec_linux_arm64 -ldflags "${LDFLAGS}" github.com/ctrsploit/ctrsploit/cmd/checksec
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o env_linux_arm64 -ldflags "${LDFLAGS}" github.com/ctrsploit/ctrsploit/cmd/env

cd -
set +x
echo "compressing binaries with upx in parallel..."
for f in bin/latest/*; do
    if [[ -f "$f" ]]; then
        (
            echo "  [upx] start  $(basename "$f")"
            upx -q "$f" >/dev/null 2>&1
            echo "  [upx] done   $(basename "$f")"
        ) &
    fi
done
wait || true
echo "compressing binaries done."
set -x

if [[ "${RELEASE_DIR}" == *release* ]]; then
    rm -f bin/release/latest
    ln -s $(realpath --relative-to=bin ${RELEASE_DIR}) bin/release/latest
fi