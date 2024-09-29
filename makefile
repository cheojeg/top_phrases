DB_URL=postgresql://root:secret@localhost:5460/top_phrases?sslmode=disable

postgres:
	docker run --name postgres_top_phrases -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5460:5432 -d postgres:12-alpine

createdb:
	docker exec -it postgres_top_phrases createdb --username=root --owner=root top_phrases

dropdb:
	docker exec -it postgres_top_phrases dropdb top_phrases

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1


migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/cheojeg/simplebank/db/sqlc Store

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock new_migration db_docs db_schema proto