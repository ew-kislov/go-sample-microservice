include .env

MODULE_NAME := $(shell go list -m)

GIT_VERSION := $(shell git describe --tags || echo none)
GIT_COMMIT := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell date +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X $(MODULE_NAME)/pkg/version.Version=$(GIT_VERSION) \
					 -X $(MODULE_NAME)/pkg/version.Commit=$(GIT_COMMIT) \
					 -X $(MODULE_NAME)/pkg/version.BuildDate=$(BUILD_DATE)

run:
	go run -ldflags "$(LDFLAGS)" ./cmd/main.go

test_unit:
	go test ./internal/...

test_integration:
	go test ./test/...

build_swagger:
	swag init -g ./internal/app.go -o ./docs

create_migration:
	migrate create -ext=sql -dir=migrations -seq init

migrate_up:
	migrate -path=migrations -database "postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path=migrations -database "postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose down

.PHONY: run test_unit test_integration create_migration migrate_up migrate_down
