include .env

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=123456 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres go-bank

dropdb:
	docker exec -it postgres12 dropdb go-bank

migration:
	migrate create -seq -ext .sql -dir ./db/migration ${name}

migrate:
	migrate -path=./db/migration -database=${DB_DSN} up

sqlc:
	sqlc generate

test:
	go test -v -cover ./...


.PHONY: migration migrate sqlc