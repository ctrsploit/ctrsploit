#!/bin/bash
#
# dqd.sh: Manages the lifecycle of a docker_archive test environment.
#
# This script is responsible for:
# 1. Starting services using docker_archive from a specified directory.
# 2. Running pre-test, main test, and post-test commands.
# 3. Ensuring the environment is torn down cleanly, regardless of test success or failure.
#

# Stop on any error, treat unset variables as errors, and propagate exit codes through pipes.
set -euo pipefail

ENV_NAME=$1
DQD_DIR=$2
PKG=$3
CMD=$4
STOP_FLAG=$5

DIR_DOCKER_ARCHIVE="/tmp/docker_archive"
DIR_SCRIPT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIR_PROJECT=$(dirname $(dirname "${DIR_SCRIPT}"))

up() {
  local dqd_dir="$1"

  if [ ! -d "${DIR_DOCKER_ARCHIVE}" ]; then
    git clone https://github.com/ssst0n3/docker_archive.git "${DIR_DOCKER_ARCHIVE}"
    ${DIR_DOCKER_ARCHIVE}/script/install_ssh_config.sh
  else
    pushd "${DIR_DOCKER_ARCHIVE}" > /dev/null
    git pull
    popd > /dev/null
  fi

  pushd "${DIR_DOCKER_ARCHIVE}/${dqd_dir}" > /dev/null
  docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
  sleep 5
  popd > /dev/null
}

down() {
  local dqd_dir="$1"
  pushd "${DIR_DOCKER_ARCHIVE}/${dqd_dir}" > /dev/null
  docker compose -f docker-compose.yml -f docker-compose.kvm.yml down
  popd > /dev/null
}

upload_test_bin() {
  local remote_host="$1"
  local pkg="$2"
  TEST_BINARY_NAME="${pkg//\//_}.test"
  pushd "${DIR_PROJECT}" > /dev/null
  scp "bin/test/${TEST_BINARY_NAME}" "${remote_host}:/ctrsploit.test"
  popd > /dev/null
}

test_cmd() {
  local remote_host="$1"
  local cmd="$2"
  ssh "$remote_host" bash -c "'
    set -euo pipefail;
    $cmd
  '"
}

do_test() {
    local remote_host="$1"
    local cmd="$2"
    local stop_flag="$3"
    if [[ -z "$stop_flag" ]]; then
      echo "No stop flag provided."
      # We explicitly redirect its standard input from /dev/null using '< /dev/null'.
      # This redirection is important because:
      #
      # 1. The outer while loop reads its input (list of files) from a process substitution.
      #    If test_cmd or any command inside it accidentally reads from standard input,
      #    it might consume the input intended for the loop and cause the loop to terminate early.
      #
      # 2. By redirecting stdin from /dev/null, we ensure that test_cmd receives an empty input,
      #    preventing it from interfering with the loop's ability to read remaining file names.
      #
      # This safeguard is necessary even if test_cmd itself does not seem to require any input.
      test_cmd "${remote_host}" "${cmd}" < /dev/null
    else
      fifo="/tmp/command_fifo.$$"
      mkfifo "$fifo"
      test_cmd "${remote_host}" "${cmd}" < /dev/null > "$fifo" &
      # stop until the stop_flag
      while IFS= read -r line; do
        echo "$line"
        if [[ "$line" == *${stop_flag}* ]]; then
          break
        fi
      done < "$fifo"
      rm $fifo
    fi
}

up "${DQD_DIR}"
upload_test_bin "${ENV_NAME}" "${PKG}"
do_test "${ENV_NAME}" "${CMD}" "${STOP_FLAG}"
down "${DQD_DIR}"