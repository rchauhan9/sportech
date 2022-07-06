#!/usr/bin/env bash

set -eu

DATABASES=(
  "football"
  "football_testing"
)

_psql() {
  psql -v ON_ERROR_STOP=1 --username "${POSTGRES_USER}" -c "$1"
}

for database in "${DATABASES[@]}"; do
  if _psql "SELECT 1 FROM pg_database WHERE datname = '$database';" | grep -q 1; then
    echo "Database ${database} already exists"
  else
    echo "Creating database ${database}"
    _psql "CREATE DATABASE $database;"
  fi
done
