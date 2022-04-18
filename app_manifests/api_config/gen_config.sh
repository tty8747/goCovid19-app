#!/usr/bin/env bash

set -o errexit
set -o pipefail

source ../env

{
printf "migration_dir: %s\n" "/app/migrations"
printf "app_endpoint: %s\n" "0.0.0.0:5000"
printf "db_host: %s\n" "${DB_ENTRYPOINT}"
printf "db_name: %s\n" "${MARIADB_DATABASE}"
printf "db_port: %s\n" "3306"
printf "db_user: %s\n" "${MARIADB_USER}"
printf "db_pass: %s\n" "${MARIADB_PASSWORD}"
printf "data_reset: %s\n" "false"
} > config.yml
