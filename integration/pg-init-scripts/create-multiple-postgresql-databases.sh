#!/bin/bash

set -e
set -u

function create_database() {
	local database=$1
	echo "  Creating user and database '$database'"
	psql -v ON_ERROR_STOP=1 -U postgres -c "CREATE DATABASE $database;"
}

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
	while ! psql -U postgres &> /dev/null ; do
		echo "Waiting for postgres to be available..."
		sleep 1
	done

	echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
	for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ',' ' '); do
		create_database $db
	done
	echo "Multiple databases created"
fi