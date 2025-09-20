# Postgres
POSTGRES_URL=postgres://root:pass.123@localhost:5432/moneytor?sslmode=disable

postgres:
	docker run --name moneytor -e POSTGRES_USER=root -e POSTGRES_PASSWORD=pass.123 -e POSTGRES_DB=moneytor -p 5432:5432 -d postgres:17-alpine

migration:
	migrate create -ext sql -dir database/migrations $(NAME)

migrateup:
	migrate -database $(POSTGRES_URL) -path database/migrations up

migrateup1:
	migrate -database $(POSTGRES_URL) -path database/migrations up 1

migratedown:
	migrate -database $(POSTGRES_URL) -path database/migrations down

migratedown1:
	migrate -database $(POSTGRES_URL) -path database/migrations down 1

sqlc:
	sqlc generate

# Tests
test:
	go test -v -count=1 ./...

# app
server:
	go run main.go
