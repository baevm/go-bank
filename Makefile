include app.env


# ======================================== DATABASE =====================================================================
postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=123456 -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres go-bank

startdb:
	docker start postgres12

dropdb:
	docker exec -it postgres12 dropdb --username=postgres go-bank

migration:
	migrate create -seq -ext .sql -dir ./db/migration ${name}

migrate-up:
	migrate -path=./db/migration -database=${DB_DSN} up

migrate-down:
	migrate -path=./db/migration -database=${DB_DSN} down

sqlc:
	sqlc generate

mock:
	mockgen -destination ./db/mock/store.go -package mockdb "go-bank/db/sqlc" Store


# ======================================== SERVER =====================================================================
test:
	go test -v -cover ./...

run:
	go run main.go

# ======================================== UTILS =====================================================================
docker-build:
	docker build -t go-bank:latest .

docker-run:
	docker run --name go-bank -p 5000:5000 -e GIN_MODE=release go-bank:latest


.PHONY: postgres createdb dropdb migration migrate-up migrate-down sqlc test run genmock