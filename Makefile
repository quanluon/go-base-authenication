include .env
export GOOSE_DRIVER=$(DRIVER)
export GOOSE_DBSTRING=$(DB_URL)
export GOOSE_MIGRATIONS_DIR=$(MIGRATIONS_DIR)

run-server:
	make sqlc
	make run

create-migration:
	goose -dir $(GOOSE_MIGRATIONS_DIR) create $(name) sql

migrate-status:
	goose -dir $(GOOSE_MIGRATIONS_DIR) status

migrate:
	goose -dir $(GOOSE_MIGRATIONS_DIR) up

rollback:
	goose -dir $(GOOSE_MIGRATIONS_DIR) down
	
sqlc:
	sqlc generate

run:
	go run main.go

test:
	go test -v ./...

