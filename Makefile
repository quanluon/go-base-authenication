include .env
export

create-migration:
	goose -dir $(GOOSE_MIGRATIONS_DIR) create $(name) sql

migrate:
	goose -dir $(GOOSE_MIGRATIONS_DIR) up

rollback:
	goose -dir $(GOOSE_MIGRATIONS_DIR) down
	
sqlc:
	sqlc generate

run:
	go run cmd/main.go

