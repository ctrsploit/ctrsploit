#!/bin/bash
#
# dqd.sh: Manages the lifecycle of a dqd test environment.
#
# This script is responsible for:
# 1. Starting services using dqd from a specified directory.
# 2. Running pre-test, main test, and post-test commands.
# 3. Ensuring the environment is torn down cleanly, regardless of test success or failure.
#

# Stop on any error, treat unset variables as errors, and propagate exit codes through pipes.
set -euo pipefail

REMOTE_HOST=$1
DQD_DIR=$2
PKG=$3
CMD=$4
STOP_FLAG=$5
START_TIMEOUT=$6
TEST_ENV_NAME=${7:-}

DIR_DQD="/tmp/dqd"
DIR_SCRIPT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIR_PROJECT=$(dirname $(dirname "${DIR_SCRIPT}"))

normalize_dqd_host() {
  local remote_host="$1"
  if [[ "${remote_host}" == dqd-* ]]; then
    echo "${remote_host}"
  else
    echo "dqd-${remote_host}"
  fi
}

REMOTE_HOST=$(normalize_dqd_host "${REMOTE_HOST}")

up() {
  local dqd_dir="$1"

  if [ ! -d "${DIR_DQD}" ]; then
    git clone https://github.com/ctrsploit/dqd.git "${DIR_DQD}"
    ${DIR_DQD}/script/install_ssh_config.sh
  else
    pushd "${DIR_DQD}" > /dev/null
    git pull
    popd > /dev/null
  fi

  pushd "${DIR_DQD}/${dqd_dir}" > /dev/null
  docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
  # until 'Reached target multi-user.target' || timeout 30
  local timeout=${START_TIMEOUT}
  local found=false
  local success_string="Reached target multi-user.target"
  local start_time=$(date +%s)
  while [ $(($(date +%s) - $start_time)) -lt ${timeout} ]; do
    if docker compose logs | sed 's/\x1b\[[0-9;]*m//g' | grep -q "${success_string}"; then
      found=true
      break
    fi
    sleep 2
  done
  if [ "${found}" = true ]; then
    echo "dqd started successfully"
  else
    echo "dqd started timeout $(( $(date +%s) - $start_time ))"
  fi
  popd > /dev/null
}

down() {
  local dqd_dir="$1"
  pushd "${DIR_DQD}/${dqd_dir}" > /dev/null
  docker compose -f docker-compose.yml -f docker-compose.kvm.yml down
  popd > /dev/null
}

cleanup() {
  local status=$?
  set +e
  down "${DQD_DIR}"
  exit "${status}"
}

upload_test_bin() {
  local remote_host="$1"
  local pkg="$2"
  TEST_BINARY_NAME="${pkg//\//_}.test"
  pushd "${DIR_PROJECT}" > /dev/null
  ssh "${remote_host}" "rm -rf /ctrsploit.test"
  scp "bin/test/${TEST_BINARY_NAME}" "${remote_host}:/ctrsploit.test"
  popd > /dev/null
}

test_cmd() {
  local remote_host="$1"
  local test_env="$2"
  local cmd="$3"
  {
    printf 'set -euo pipefail\n'
    printf 'export TEST_ENV=%q\n' "${test_env}"
    printf '%s\n' "${cmd}"
  } | ssh "$remote_host" bash -s
}

do_test() {
    local remote_host="$1"
    local test_env="$2"
    local cmd="$3"
    local stop_flag="$4"
    if [[ -z "$stop_flag" || "$stop_flag" == "null" ]]; then
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
      test_cmd "${remote_host}" "${test_env}" "${cmd}" < /dev/null
    else
      fifo="/tmp/command_fifo.$$"
      mkfifo "$fifo"
      test_cmd "${remote_host}" "${test_env}" "${cmd}" < /dev/null > "$fifo" &
      test_pid=$!
      matched=false
      # stop until the stop_flag
      while IFS= read -r line; do
        echo "$line"
        if [[ "$line" == *${stop_flag}* ]]; then
          matched=true
          break
        fi
      done < "$fifo"
      rm -f "$fifo"
      if [[ "${matched}" == true ]]; then
        # The remote command printed its success line: the guest may be wedged
        # from there on (a parked D-state worker can take sshd with it), so a
        # dirty session exit does not invalidate the result that was already
        # printed and verified. Stop the session and report success.
        kill "${test_pid}" 2>/dev/null || true
        wait "${test_pid}" 2>/dev/null || true
        return 0
      fi
      # The flag never appeared — this is a real failure. Return the remote
      # command's status: the reader loop ends at EOF when the session dies
      # (timeout, guest wedge) and would otherwise report rm's exit status 0,
      # turning a failure into a false green.
      local rc=0
      wait "${test_pid}" || rc=$?
      return "${rc}"
    fi
}

up "${DQD_DIR}"
trap cleanup EXIT
upload_test_bin "${REMOTE_HOST}" "${PKG}"
do_test "${REMOTE_HOST}" "${TEST_ENV_NAME}" "${CMD}" "${STOP_FLAG}"
