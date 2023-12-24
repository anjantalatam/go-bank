include .env
export

postgres:
	docker run --name postgres16 -e POSTGRES_USER=$(POSTGRES_USER) -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -p 5432:5432 -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) simple_bank

dropdb:
	docker exec -it postgres16 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DATABASE_URL)?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "$(DATABASE_URL)?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown