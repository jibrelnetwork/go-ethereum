#!/bin/sh -e

echo "Go-Ethereum Fork Information:"
geth version

external_ip_option=""
if [ ! -z ${GETIP} ]; then
    EXTERNAL_IP=$(wget -q -O - https://api.ipify.org)
	if expr "$EXTERNAL_IP" : '[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}$' >/dev/null; then 
		echo "Using '${EXTERNAL_IP}' as an external IP address..."
		external_ip_option="--nat extip:${EXTERNAL_IP}"
	else
		echo "Got wrong external IP address: '${EXTERNAL_IP}'. Ignoring."
	fi
fi

if [ -z "${POSTGRES_DB_HOST}" ] || [ -z "${POSTGRES_DB_USER}" ] || [ -z "${POSTGRES_DB_PASS}" ]; then
    extdb_option=""
    echo "Database information is not set in env, migrations will be skipped. You can still pass extdb option manually."
else
    DB_PORT=${POSTGRES_DB_PORT:-5432}
    DB_NAME=${POSTGRES_DB_NAME:-jsearch-raw}
    DB_SSL_MODE=${POSTGRES_DB_SSL_MODE:-disable}

    extdb_option="-extdb=postgres://${POSTGRES_DB_USER}:${POSTGRES_DB_PASS}@${POSTGRES_DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"

    echo "Executing migrations..."
    goose -dir /schema_migrations/ postgres "host=${POSTGRES_DB_HOST} port=${DB_PORT} user=${POSTGRES_DB_USER} password=${POSTGRES_DB_PASS} dbname=${DB_NAME} sslmode=${DB_SSL_MODE}" up

    echo "Updating permissions..."
    for pg_readonly_user in ${POSTGRES_READONLY_USERS//,/ } ; do
        echo " ... for user ${pg_readonly_user}"
        PGPASSWORD="${POSTGRES_DB_PASS}" psql -h "${POSTGRES_DB_HOST}" -p ${DB_PORT} -U "${POSTGRES_DB_USER}" -d "${DB_NAME}" \
            -c "GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO \"${pg_readonly_user}\";" \
            -c "GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"${pg_readonly_user}\";"
    done
fi

echo "Starting geth..."
exec geth ${extdb_option} ${external_ip_option} "$@"
