postgres:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:17-alpine

createdb:
	docker exec -it postgres17 createdb -U postgres -O postgres bank

dropdb:
	docker exec -it postgres17 dropdb -U postgres bank

migrateup:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test