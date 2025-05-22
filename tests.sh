#!/bin/bash
set -e

CONTAINER_NAME="test_postgr"
DB_PORT="5437"
first_time=false

# 1. Запуск/переиспользование контейнера
if ! docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}\$"; then
    echo "Starting PostgreSQL container..."
    first_time=true
    docker run -d \
        -p "${DB_PORT}:5432" \
        -e POSTGRES_USER=test \
        -e POSTGRES_PASSWORD=test \
        -e POSTGRES_DB=test \
        --name "$CONTAINER_NAME" \
        postgres:16
else
    echo "Restarting existing container..."
    docker start "$CONTAINER_NAME"
fi


export TEST_DB_URL="postgres://test:test@localhost:${DB_PORT}/test?sslmode=disable"

if "$first_time"; then
    go build -tags=tests -o tests ./cmd/tests/main.go && ./tests
fi

echo "Running tests with DB URL: $TEST_DB_URL"

# Получить список пакетов с интеграционными тестами
PKGS=$(go list -tags=integration -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

# Запустить тесты только для них
go test -v -tags=integration $PKGS
# go test -v -tags=integration ./...

docker stop "$CONTAINER_NAME"