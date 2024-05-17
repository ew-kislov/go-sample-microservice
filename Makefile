include app.env

create_migration:
	migrate create -ext=sql -dir=migrations -seq init

migrate_up:
	migrate -path=migrations -database "postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose up

migrate_down:
	migrate -path=migrations -database "postgresql://${DATABASE_USER}:${DATABASE_PASSWORD}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable" -verbose down

.PHONY: create_migration migrate_up migrate_down
