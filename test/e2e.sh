#!/bin/bash
set -euo pipefail

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
  local docker_archive_dir="/tmp/docker_archive"

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

test_cmd() {
  local remote_host="$1"
  local cmd="$2"

  # TODO: remove this when the dev image is available in the remote registry
  docker save ghcr.io/ctrsploit/ctrsploit-dev > /tmp/dev-image.tar
  scp /tmp/dev-image.tar "${remote_host}:/root/dev-image.tar"

  ssh "$remote_host" bash -c "'
    set -euo pipefail
    mkdir -p ctrsploit
    tar xzf ctrsploit.tar.gz -C ctrsploit
    cd ctrsploit
    git config --global --add safe.directory /root/ctrsploit
    # TODO: remove this when the dev image is available in the remote registry
    docker load < /root/dev-image.tar
    $cmd
  '"
}

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
    ENV_CMD=$(yq eval ".test_envs[$i].cmd" "$e2e_file")
    DQD_DIR=$(yq eval ".test_envs[$i].dqd_dir" "$e2e_file")

    echo "----------------------------------"
    echo "TEST_ENV = ${ENV_NAME}"
    echo "TEST_CMD = ${ENV_CMD}"
    echo "DQD_DIR  = ${DQD_DIR}"

    startup_testEnv "${DQD_DIR}"
    upload_codebase "${ENV_NAME}"
    test_cmd "${ENV_NAME}" "${ENV_CMD}"
  done
done < <(find . -type f -name "e2e.yml")
