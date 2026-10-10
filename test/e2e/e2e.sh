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
    # retries: how many times to re-run a failed env (each dqd retry boots a
    # FRESH VM). Kernel-exploit e2es are stochastic per boot — a reclaim race
    # or a panic wedge can fail one boot's conditions while the next boot's
    # SLUB state succeeds — so vulnerable envs opt in with retries: 2.
    RETRIES=$(yq eval ".test_envs[$i].retries // .defaults.retries // 1" "$e2e_file")

    if [[ -z "$REMOTE_HOST" || "$REMOTE_HOST" == "null" ]]; then
        REMOTE_HOST="$ENV_NAME"
    fi

    if [[ -z "$START_TIMEOUT" || "$START_TIMEOUT" == "null" ]]; then
        START_TIMEOUT=60
    fi

    if ! [[ "$RETRIES" =~ ^[0-9]+$ ]] || [[ "$RETRIES" -lt 1 ]]; then
        RETRIES=1
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
    echo "RETRIES = ${RETRIES}"

    make test.bin PKG=${PKG}

    ATTEMPT=1
    until [[ "${ATTEMPT}" -gt "${RETRIES}" ]]; do
      if [[ "${ENV_KIND}" == "plain" ]]; then
        if (
          export TEST_ENV="${ENV_NAME}"
          eval "${CMD}"
        ); then break; fi
      elif [[ "${ENV_KIND}" == "dqd" ]]; then
        if ${DIR_SCRIPT}/dqd.sh "${REMOTE_HOST}" "${DQD_DIR}" "${PKG}" "${CMD}" "${STOP_FLAG}" "${START_TIMEOUT}" "${ENV_NAME}"; then break; fi
      else
        break
      fi
      echo "e2e attempt ${ATTEMPT}/${RETRIES} for ${ENV_NAME} failed"
      ATTEMPT=$((ATTEMPT + 1))
    done

    if [[ "${ATTEMPT}" -gt "${RETRIES}" ]]; then
      echo "e2e ${ENV_NAME} failed after ${RETRIES} attempt(s)" >&2
      exit 1
    fi
  done
done < <(find ${DIR_SEARCH} -type f -name "e2e.yml")
popd > /dev/null

if [[ -n "${TARGET_TEST_ENV}" && "${MATCHED_TEST_ENV}" != true ]]; then
  echo "No e2e test environment matched TEST_ENV=${TARGET_TEST_ENV}" >&2
  exit 1
fi
