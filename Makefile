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
	$(MAKE) mock

mock:
	mockgen -destination ./db/mock/store.go -package mockdb "go-bank/db/sqlc" Store


# ======================================== SERVER =====================================================================
test:
	go test -v -cover ./...

run:
	go run main.go

docker-build:
	docker build -f deploy/docker/Dockerfile -t go-bank:latest .

docker-run:
	docker run --name go-bank -p 5000:5000 -e GIN_MODE=release go-bank:latest

docker-compose:
	docker-compose -f deploy/docker/docker-compose.yml up 

# ======================================== GRPC ======================================================================
grpc-compile:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=go-bank \
	proto/*.proto

# ======================================== UTILS ======================================================================


.PHONY: postgres createdb dropdb migration migrate-up migrate-down sqlc test run genmock grpc-compile