#!/bin/bash
set -ex

# Global variables
SCRIPT_DIR=$(dirname "$(readlink -f "${BASH_SOURCE[0]}")")
PROJECT_DIR=$(dirname ${SCRIPT_DIR})
RELEASE_DIR=""
VERSION=""

# Setup release directory and symbolic links
setup_release_dir() {
    # Load version information
    source ${SCRIPT_DIR}/version.sh
    
    cd ${PROJECT_DIR}
    mkdir -p ${RELEASE_DIR}
    
    # Create bin/latest symbolic link
    rm -f bin/latest
    ln -s $(realpath --relative-to=bin ${RELEASE_DIR}) bin/latest
    
    cd ${RELEASE_DIR}
}

# Build a single binary file
# Parameters:
#   $1: Binary file name (without extension)
#   $2: GOOS (e.g., linux)
#   $3: GOARCH (e.g., amd64)
#   $4: Package path (e.g., github.com/ctrsploit/ctrsploit/cmd/ctrsploit)
build_binary() {
    local binary_name=$1
    local goos=$2
    local goarch=$3
    local package_path=$4
    local output_name="${binary_name}_${goos}_${goarch}"
    
    CGO_ENABLED=0 GOOS=${goos} GOARCH=${goarch} go build \
        -o ${output_name} \
        -ldflags "${LDFLAGS}" \
        ${package_path}
}

# Build binaries for all platforms
build_all_binaries() {
    local commands=("ctrsploit" "checksec" "env")
    local platforms=("linux:amd64" "linux:arm64")
    
    for platform in "${platforms[@]}"; do
        IFS=':' read -r goos goarch <<< "${platform}"
        
        for cmd in "${commands[@]}"; do
            local package_path="github.com/ctrsploit/ctrsploit/cmd/${cmd}"
            build_binary "${cmd}" "${goos}" "${goarch}" "${package_path}"
        done
    done
}

# Compress binaries using upx
compress_binaries() {
    cd ${PROJECT_DIR}
    set +x
    
    # Choose upx compression level:
    # - release versions (bin/release/*): highest compression (--best)
    # - other versions (bin/dev/*, dirty, etc.): fastest compression (--fast)
    local upx_args="-q --fast"
    if [[ "${RELEASE_DIR}" == bin/release/* ]]; then
        echo "use upx highest compression for release binaries (--best)"
        upx_args="-q --best"
    else
        echo "use upx fastest compression for non-release binaries (--fast)"
    fi
    
    echo "compressing binaries with upx in parallel..."
    for f in bin/latest/*; do
        if [[ -f "$f" ]]; then
            (
                echo "  [upx] start  $(basename "$f")"
                upx ${upx_args} "$f" >/dev/null 2>&1
                echo "  [upx] done   $(basename "$f")"
            ) &
        fi
    done
    
    wait || true
    echo "compressing binaries done."
    
    set -x
}

# Create release/latest link if it's a release version
create_release_link() {
    if [[ "${RELEASE_DIR}" == *release* ]]; then
        rm -f bin/release/latest
        ln -s $(realpath --relative-to=bin ${RELEASE_DIR}) bin/release/latest
    fi
}

setup_release_dir
build_all_binaries
compress_binaries
create_release_link
