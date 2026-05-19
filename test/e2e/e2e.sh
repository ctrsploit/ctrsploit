#!/bin/bash
set -euo pipefail

DIR_SEARCH=${1:-.}
TARGET_TEST_ENV=${TEST_ENV:-}
unset TEST_ENV
# Get the absolute path of the script directory
DIR_SCRIPT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIR_PROJECT=$(dirname $(dirname "${DIR_SCRIPT}"))
MATCHED_TEST_ENV=false

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
    # Extract the test environment details, falling back to defaults.
    ENV_NAME=$(yq eval ".test_envs[$i].name" "$e2e_file")
    if [[ -n "${TARGET_TEST_ENV}" && "${ENV_NAME}" != "${TARGET_TEST_ENV}" ]]; then
      continue
    fi
    MATCHED_TEST_ENV=true

    REMOTE_HOST=$(yq eval ".test_envs[$i].remote_host // .defaults.remote_host" "$e2e_file")
    ENV_KIND=$(yq eval ".test_envs[$i].kind // .defaults.kind" "$e2e_file")
    DQD_DIR=$(yq eval ".test_envs[$i].dqd_dir // .defaults.dqd_dir" "$e2e_file")
    PKG=$(yq eval ".test_envs[$i].pkg // .defaults.pkg" "$e2e_file")
    CMD=$(yq eval ".test_envs[$i].cmd // .defaults.cmd" "$e2e_file")
    STOP_FLAG=$(yq eval ".test_envs[$i].stop_flag // .defaults.stop_flag" "$e2e_file")
    START_TIMEOUT=$(yq eval ".test_envs[$i].start_timeout // .defaults.start_timeout" "$e2e_file")

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
      (
        export TEST_ENV="${ENV_NAME}"
        eval "${CMD}"
      )
    elif [[ "${ENV_KIND}" == "dqd" ]]; then
      ${DIR_SCRIPT}/dqd.sh "${REMOTE_HOST}" "${DQD_DIR}" "${PKG}" "${CMD}" "${STOP_FLAG}" "${START_TIMEOUT}" "${ENV_NAME}"
    fi
  done
done < <(find ${DIR_SEARCH} -type f -name "e2e.yml")
popd > /dev/null

if [[ -n "${TARGET_TEST_ENV}" && "${MATCHED_TEST_ENV}" != true ]]; then
  echo "No e2e test environment matched TEST_ENV=${TARGET_TEST_ENV}" >&2
  exit 1
fi
