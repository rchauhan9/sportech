#!/usr/bin/env bash

set -Eeo pipefail

mkdir -p /startup-initdb.d/

entrypoint="$(which docker-entrypoint.sh)"
source "$entrypoint"

docker_setup_env
docker_create_db_directories

# restart script as postgres when running as root
if [ "$(id -u)" = '0' ]; then
   exec gosu postgres "$BASH_SOURCE" "$@"
fi

if [ -n "$DATABASE_ALREADY_EXISTS" ]; then
  docker_temp_server_start
  docker_process_init_files /startup-initdb.d/*
  docker_temp_server_stop
fi

exec "$entrypoint" "postgres"
