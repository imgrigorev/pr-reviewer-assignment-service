#!/bin/bash

source .env

export MIGRATION_DSN="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@db:5432/$POSTGRES_DB?sslmode=disable"

echo "Waiting for database to be ready..."
sleep 5

echo "Running migrations with DSN: postgres://$POSTGRES_USER:****@db:5432/$POSTGRES_DB?sslmode=disable"

migrate -path "${MIGRATION_DIR}" -database "${MIGRATION_DSN}" up

if [ $? -eq 0 ]; then
    echo "Migrations completed successfully"
else
    echo "Migrations failed"
    exit 1
fi
