#!/usr/bin/env bash
set -euo pipefail

DQD_ROOT=${DQD_ROOT:-/tmp/dqd}
DQD_REPO=${DQD_REPO:-https://github.com/ctrsploit/dqd.git}
DQD_SUCCESS_STRING=${DQD_SUCCESS_STRING:-Reached target multi-user.target}

usage() {
  cat <<'USAGE'
Usage:
  dqd-lab.sh list [--no-update]
  dqd-lab.sh up <dqd_dir> [--timeout seconds] [--no-update] [--no-install-ssh-config]
  dqd-lab.sh down <dqd_dir>
  dqd-lab.sh ps <dqd_dir>
  dqd-lab.sh logs <dqd_dir> [docker compose logs args...]
  dqd-lab.sh ssh <remote_host> [-- command]

Environment:
  DQD_ROOT             dqd checkout path, default /tmp/dqd
  DQD_REPO             dqd git repository, default https://github.com/ctrsploit/dqd.git
  DQD_SUCCESS_STRING   boot log marker, default "Reached target multi-user.target"
USAGE
}

die() {
  printf 'error: %s\n' "$*" >&2
  exit 1
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "missing required command: $1"
}

require_docker_compose() {
  need_cmd docker
  docker compose version >/dev/null 2>&1 || die "docker compose is unavailable"
}

validate_dqd_dir() {
  local dqd_dir=$1
  [[ -n "$dqd_dir" ]] || die "dqd_dir is required"
  [[ "$dqd_dir" != /* ]] || die "dqd_dir must be relative to DQD_ROOT"
  [[ "$dqd_dir" != *".."* ]] || die "dqd_dir must not contain '..'"
}

dqd_path() {
  local dqd_dir=$1
  validate_dqd_dir "$dqd_dir"
  printf '%s/%s\n' "$DQD_ROOT" "$dqd_dir"
}

ensure_dqd() {
  local update=${1:-1}
  local install_ssh_config=${2:-1}

  need_cmd git
  if [[ ! -d "$DQD_ROOT/.git" ]]; then
    git clone "$DQD_REPO" "$DQD_ROOT"
  elif [[ "$update" == 1 ]]; then
    git -C "$DQD_ROOT" pull --ff-only
  fi

  if [[ "$install_ssh_config" == 1 && -x "$DQD_ROOT/script/install_ssh_config.sh" ]]; then
    "$DQD_ROOT/script/install_ssh_config.sh"
  fi
}

compose_in() {
  local dir=$1
  shift
  [[ -d "$dir" ]] || die "dqd directory does not exist: $dir"
  [[ -f "$dir/docker-compose.yml" ]] || die "missing docker-compose.yml in $dir"
  [[ -f "$dir/docker-compose.kvm.yml" ]] || die "missing docker-compose.kvm.yml in $dir"
  (cd "$dir" && docker compose -f docker-compose.yml -f docker-compose.kvm.yml "$@")
}

wait_for_boot() {
  local dir=$1
  local timeout=$2
  local start now

  start=$(date +%s)
  while true; do
    if compose_in "$dir" logs --no-color 2>/dev/null | grep -Fq "$DQD_SUCCESS_STRING"; then
      printf 'dqd started successfully: %s\n' "$dir"
      return 0
    fi

    now=$(date +%s)
    if (( now - start >= timeout )); then
      printf 'warning: dqd start timed out after %ss: %s\n' "$((now - start))" "$dir" >&2
      return 1
    fi
    sleep 2
  done
}

normalize_host() {
  local host=$1
  if [[ "$host" == dqd-* ]]; then
    printf '%s\n' "$host"
  else
    printf 'dqd-%s\n' "$host"
  fi
}

cmd_list() {
  local update=1
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --no-update) update=0; shift ;;
      -h|--help) usage; exit 0 ;;
      *) die "unknown list option: $1" ;;
    esac
  done

  ensure_dqd "$update" 0
  find "$DQD_ROOT" -mindepth 2 -maxdepth 4 -name docker-compose.yml -print \
    | sed "s#^$DQD_ROOT/##; s#/docker-compose.yml\$##" \
    | sort
}

cmd_up() {
  local dqd_dir=${1:-}
  local timeout=60
  local update=1
  local install_ssh_config=1
  [[ -n "$dqd_dir" ]] || die "up requires <dqd_dir>"
  shift

  while [[ $# -gt 0 ]]; do
    case "$1" in
      --timeout) timeout=${2:?--timeout requires seconds}; shift 2 ;;
      --timeout=*) timeout=${1#*=}; shift ;;
      --no-update) update=0; shift ;;
      --no-install-ssh-config) install_ssh_config=0; shift ;;
      -h|--help) usage; exit 0 ;;
      *) die "unknown up option: $1" ;;
    esac
  done

  [[ "$timeout" =~ ^[0-9]+$ ]] || die "timeout must be an integer"
  ensure_dqd "$update" "$install_ssh_config"
  local dir
  dir=$(dqd_path "$dqd_dir")
  require_docker_compose
  compose_in "$dir" up -d
  wait_for_boot "$dir" "$timeout" || return 0
}

cmd_down() {
  local dqd_dir=${1:-}
  [[ -n "$dqd_dir" ]] || die "down requires <dqd_dir>"
  local dir
  dir=$(dqd_path "$dqd_dir")
  require_docker_compose
  compose_in "$dir" down
}

cmd_ps() {
  local dqd_dir=${1:-}
  [[ -n "$dqd_dir" ]] || die "ps requires <dqd_dir>"
  local dir
  dir=$(dqd_path "$dqd_dir")
  require_docker_compose
  compose_in "$dir" ps
}

cmd_logs() {
  local dqd_dir=${1:-}
  [[ -n "$dqd_dir" ]] || die "logs requires <dqd_dir>"
  shift
  local dir
  dir=$(dqd_path "$dqd_dir")
  require_docker_compose
  if [[ $# -eq 0 ]]; then
    compose_in "$dir" logs --tail=200
  else
    compose_in "$dir" logs "$@"
  fi
}

cmd_ssh() {
  local host=${1:-}
  [[ -n "$host" ]] || die "ssh requires <remote_host>"
  shift
  if [[ "${1:-}" == "--" ]]; then
    shift
  fi
  host=$(normalize_host "$host")
  need_cmd ssh
  if [[ $# -eq 0 ]]; then
    ssh "$host"
  else
    ssh "$host" "$@"
  fi
}

main() {
  local cmd=${1:-}
  case "$cmd" in
    list) shift; cmd_list "$@" ;;
    up) shift; cmd_up "$@" ;;
    down) shift; cmd_down "$@" ;;
    ps) shift; cmd_ps "$@" ;;
    logs) shift; cmd_logs "$@" ;;
    ssh) shift; cmd_ssh "$@" ;;
    -h|--help|help|"") usage ;;
    *) die "unknown command: $cmd" ;;
  esac
}

main "$@"
