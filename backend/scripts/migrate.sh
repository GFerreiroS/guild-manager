#!/bin/sh
set -e

attempt=1
max_retries=${MIGRATION_MAX_RETRIES:-3}

until [ $attempt -gt $max_retries ]
do
    echo "Running migrations (attempt $attempt/$max_retries)"
    go run cmd/migrate/main.go -action up && break
    
    attempt=$((attempt+1))
    sleep 5
    
    if [ $attempt -gt $max_retries ]; then
        echo "Migration failed after $max_retries attempts"
        exit 1
    fi
done