include .env

MODULE_NAME := $(shell go list -m)

GIT_VERSION := $(shell git describe --tags || echo none)
GIT_COMMIT := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X $(MODULE_NAME)/pkg.Version=$(GIT_VERSION) \
					 -X $(MODULE_NAME)/pkg.Commit=$(GIT_COMMIT) \
					 -X $(MODULE_NAME)/pkg.BuildDate=$(BUILD_DATE)

run:
	go run -ldflags "$(LDFLAGS)" ./cmd/main.go

tests_unit:
	go test -v $(shell go list ./... | grep -v /test)

tests_integration:
	go test -v ./internal/test

create_migration:
	migrate create -ext=sql -dir=migrations -seq init

migrate_up:
	migrate -path=migrations -database "postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path=migrations -database "postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose down

.PHONY: create_migration migrate_up migrate_down
