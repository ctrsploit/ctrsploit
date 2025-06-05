#!/bin/bash
set -euo pipefail

docker_archive_dir="/tmp/docker_archive"

# Check for required commands.
for cmd in yq git docker scp ssh; do
  if ! command -v "$cmd" > /dev/null 2>&1; then
    echo "Error: $cmd is required but not installed. Aborting."
    exit 1
  fi
done

upload_codebase() {
  local remote_host="$1"
  local tmp_tar="/tmp/ctrsploit.tar.gz"
  # Determine the directory containing this script.
  local script_dir
  script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
  # The project directory is two levels up from the script directory.
  local project_dir
  project_dir=$(dirname "$script_dir")

  pushd "$project_dir" > /dev/null
  tar czf "$tmp_tar" --exclude='.idea' --exclude='bin' .
  popd > /dev/null

  scp "$tmp_tar" "${remote_host}:"
}

startup_testEnv() {
  local dqd_dir="$1"

  if [ ! -d "$docker_archive_dir" ]; then
    git clone https://github.com/ssst0n3/docker_archive.git "$docker_archive_dir"
    ${docker_archive_dir}/script/install_ssh_config.sh
  else
    pushd "$docker_archive_dir" > /dev/null
    git pull
    popd > /dev/null
  fi

  pushd "${docker_archive_dir}/${dqd_dir}" > /dev/null
  docker compose -f docker-compose.yml -f docker-compose.kvm.yml up -d
  sleep 1
  popd > /dev/null
}

stop_testEnv() {
  local dqd_dir="$1"
  pushd "${docker_archive_dir}/${dqd_dir}" > /dev/null
  docker compose -f docker-compose.yml -f docker-compose.kvm.yml down
  popd > /dev/null
}

test_cmd() {
  local remote_host="$1"
  local pre_cmd="$2"
  local cmd="$3"

  # TODO: remove this when the dev image is available in the remote registry
  docker save ghcr.io/ctrsploit/ctrsploit-dev > /tmp/dev-image.tar
  scp /tmp/dev-image.tar "${remote_host}:/root/dev-image.tar"
  ssh "$remote_host" bash -c "'
    set -euo pipefail;
    $pre_cmd
    $cmd
  '"
}

do_test() {
    local remote_host="$1"
    local pre_cmd="$2"
    local cmd="$3"
    local stop_flag="$4"
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
      test_cmd "${remote_host}" "${pre_cmd}" "${cmd}" < /dev/null
    else
      fifo="/tmp/command_fifo.$$"
      mkfifo "$fifo"
      test_cmd "${remote_host}" "${pre_cmd}" "${cmd}" < /dev/null > "$fifo" &
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

SEARCH_DIR=${1:-.}

# Loop through all e2e.yml files in the current directory recursively.
while IFS= read -r e2e_file; do
  echo "=================================="
  echo "Handling: ${e2e_file}"
  num_envs=$(yq eval '.test_envs | length' "$e2e_file")

  if [ "$num_envs" -eq 0 ]; then
    echo "No test environments found in ${e2e_file}"
    continue
  fi

  for (( i = 0; i < num_envs; i++ )); do
    # Extract the test environment details.
    ENV_NAME=$(yq eval ".test_envs[$i].name" "$e2e_file")
    ENV_PRE_CMD=$(yq eval ".test_envs[$i].pre_cmd" "$e2e_file")
    ENV_CMD=$(yq eval ".test_envs[$i].cmd" "$e2e_file")
    DQD_DIR=$(yq eval ".test_envs[$i].dqd_dir" "$e2e_file")
    STOP_FLAG=$(yq eval ".test_envs[$i].stop_flag" "$e2e_file")

    echo "----------------------------------"
    echo "TEST_ENV = ${ENV_NAME}"
    echo "PRE_CMD = ${ENV_PRE_CMD}"
    echo "TEST_CMD = ${ENV_CMD}"
    echo "DQD_DIR  = ${DQD_DIR}"
    echo "STOP_FLAG  = ${STOP_FLAG}"

    startup_testEnv "${DQD_DIR}"
    upload_codebase "${ENV_NAME}"
    do_test "${ENV_NAME}" "${ENV_PRE_CMD}" "${ENV_CMD}" "${STOP_FLAG}"
    stop_testEnv "${DQD_DIR}"
  done
done < <(find ${SEARCH_DIR} -type f -name "e2e.yml")
