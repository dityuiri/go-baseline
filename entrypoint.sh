#!/bin/sh
set -e

#./migrate -database "$DB_DRIVER://$DB_MIGRATION_USER:$DB_MIGRATION_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable&search_path=$DB_SCHEMA" -path db/migrations up

exec "$@"
