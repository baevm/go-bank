#!/bin/sh

source "app.env"

set -e

echo "run db migration"
/app/go-bank/migrate -path /app/go-bank/migration -database "$DB_DSN" -verbose up

echo "start the app"
exec "$@"