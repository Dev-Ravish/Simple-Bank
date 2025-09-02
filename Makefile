postgres: 
	docker run -p 8800:5432 --name postgres12 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb: 
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb --username=root simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:8800/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:8800/simple_bank?sslmode=disable" -verbose drop down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...
	
.PHONY: postgres createdb dropdb