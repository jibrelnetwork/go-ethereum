#!/bin/sh -e

echo "Go-Ethereum Fork Information:"
geth version

if [ -z "${POSTGRES_DB_HOST}" ] || [ -z "${POSTGRES_DB_USER}" ] || [ -z "${POSTGRES_DB_PASS}" ]; then
    extdb_option=""
    echo "Database information is not set in env, migrations will be skipped. You can still pass extdb option manually."
else
    DB_PORT=${POSTGRES_DB_PORT:-5432}
    DB_NAME=${POSTGRES_DB_NAME:-jsearch-raw}
    DB_SSL_MODE=${POSTGRES_DB_SSL_MODE:-disable}

    extdb_option="-extdb=postgres://${POSTGRES_DB_USER}:${POSTGRES_DB_PASS}@${POSTGRES_DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

    echo "Executing migrations..."
    goose -dir /migrations/ postgres "host=${POSTGRES_DB_HOST} port=${DB_PORT} user=${POSTGRES_DB_USER} password=${POSTGRES_DB_PASS} dbname=${DB_NAME} sslmode=${DB_SSL_MODE}" up
fi

echo "Starting geth..."
exec geth ${extdb_option} "$@"
