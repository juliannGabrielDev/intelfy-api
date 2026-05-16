.PHONY: db-up db-down migrate-up migrate-down migrate-reset migrate-create migrate-version sqlc-generate

DB_CONTAINER ?= intelfy_db
DB_URL ?= postgres://postgres:pass123@localhost:5445/intelfy_db?sslmode=disable
MIGRATIONS_DIR = db/migrations

db-up:
	docker compose up -d db

db-down:
	docker compose down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-reset:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" drop -f
	$(MAKE) migrate-up

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

sqlc-generate:
	sqlc generate
