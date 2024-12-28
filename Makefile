include .env
export

postgres:
	docker run --name whale-users-postgres -p 5432:5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PWORD) -d postgres:12-alpine

createdb:
	docker exec -it whale-users-postgres createdb --username=$(DB_USER) --owner=admin whale_wake_users

dropdb:
	docker exec -it whale-users-postgres dropdb -U $(DB_USER) whale_wake_users

migrateup:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PWORD)@localhost:5432/whale_wake_users?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://$(DB_USER):$(DB_PWORD)@localhost:5432/whale_wake_users?sslmode=disable" -verbose down


.PHONY: postgres createdb dropdb migrateup migratedown