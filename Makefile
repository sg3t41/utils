.PHONY: migration migrate-up migrate-down migrate-status migrate-reset migrate-force

DATABASE_URL := postgres://postgres:password@localhost:5432/apidb?sslmode=disable
MIGRATIONS_PATH := ./migrations

migration:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $$name

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

migrate-status:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" version

migrate-reset:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" drop -f
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $$version

install-migrate-tool:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest