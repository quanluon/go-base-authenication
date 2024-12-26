include .env
export GOOSE_DRIVER=$(GOOSE_DRIVER)
export GOOSE_DBSTRING=$(GOOSE_DBSTRING)
export GOOSE_MIGRATIONS_DIR=$(GOOSE_MIGRATIONS_DIR)

run-server:
	make sqlc
	make run

create-migration:
	goose -dir $(GOOSE_MIGRATIONS_DIR) create $(name) sql

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

