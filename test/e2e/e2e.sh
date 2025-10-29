#!/bin/bash
set -euo pipefail

DIR_SEARCH=${1:-.}
# Get the absolute path of the script directory
DIR_SCRIPT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIR_PROJECT=$(dirname $(dirname "${DIR_SCRIPT}"))

pushd "${DIR_PROJECT}" > /dev/null

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
    REMOTE_HOST=$(yq eval ".test_envs[$i].remote_host" "$e2e_file")
    ENV_KIND=$(yq eval ".test_envs[$i].kind" "$e2e_file")
    DQD_DIR=$(yq eval ".test_envs[$i].dqd_dir" "$e2e_file")
    PKG=$(yq eval ".test_envs[$i].pkg" "$e2e_file")
    CMD=$(yq eval ".test_envs[$i].cmd" "$e2e_file")
    STOP_FLAG=$(yq eval ".test_envs[$i].stop_flag" "$e2e_file")
    START_TIMEOUT=$(yq eval ".test_envs[$i].start_timeout" "$e2e_file")

    if [[ -z "$REMOTE_HOST" || "$REMOTE_HOST" == "null" ]]; then
        REMOTE_HOST="$ENV_NAME"
    fi

    if [[ -z "$START_TIMEOUT" || "$START_TIMEOUT" == "null" ]]; then
        START_TIMEOUT=60
    fi

    echo "----------------------------------"
    echo "TEST_ENV = ${ENV_NAME}"
    echo "REMOTE_HOST = ${REMOTE_HOST}"
    echo "ENV_KIND = ${ENV_KIND}"
    echo "DQD_DIR = ${DQD_DIR}"
    echo "PKG = ${PKG}"
    echo "TEST_CMD = ${CMD}"
    echo "STOP_FLAG  = ${STOP_FLAG}"
    echo "START_TIMEOUT = ${START_TIMEOUT}"

    make test.bin PKG=${PKG}

    if [[ "${ENV_KIND}" == "plain" ]]; then
      eval "${CMD}"
    elif [[ "${ENV_KIND}" == "dqd" ]]; then
      ${DIR_SCRIPT}/dqd.sh "${REMOTE_HOST}" "${DQD_DIR}" "${PKG}" "${CMD}" "${STOP_FLAG}" "${START_TIMEOUT}"
    fi
  done
done < <(find ${DIR_SEARCH} -type f -name "e2e.yml")
popd > /dev/null
