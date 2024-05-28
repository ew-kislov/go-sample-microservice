# Go Sample Microservice

## Description

This is a sample microservice written in Go considering best practices.
It can be used as a template.

## Features

- REST API endpoints
- JWT authorization
- SQL database as a storage
- Database migrations
- Auto generated OpenApi documentation
- End-to-End logging (with Request ID). Includes HTTP requests logging, database queries logging
- Testing infrastructure (mocks, integration tests)
- Linter

## Prerequisites

- `migrate` (https://github.com/golang-migrate/migrate)
- `golangci-lint` (https://github.com/golangci/golangci-lint)
- `swag` (https://github.com/swaggo/swag)

## How to run app

1. Create database PostgreSQL
2. Run migrations: `make migrate_up`
3. Install dependencies: `go mod download`
4. Run app: `make run`
