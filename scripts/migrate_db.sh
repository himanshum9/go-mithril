#!/bin/bash
set -e

# Load environment variables
source "$(dirname "$0")/.env"


echo "Running database migrations..."
migrate -path ./migrations -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" up

echo "Seeding initial data..."
# psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f ./seeds/initial_data.sql

echo "Database migration and seeding completed."