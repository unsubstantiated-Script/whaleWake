include .env
export

postgres:
	docker run --name whale-users-postgres -p 5432:5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PWORD) -d postgres:12-alpine

createdb:
	docker exec -it whale-users-postgres createdb --username=$(DB_USER) --owner=admin whale_wake_users

dropdb:
	docker exec -it whale-users-postgres dropdb -U $(DB_USER) whale_wake_users

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

testmigrateup:
	migrate -path db/migration -database "postgresql://testuser:testpassword@localhost:5432/whale_wake_users?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup testmigrateup migratedown sqlc test server