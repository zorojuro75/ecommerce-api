.PHONY: run build down logs restart clean test

run:
	go run cmd/api/main.go

docker-build:
	docker-compose build --no-cache

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f api

docker-restart:
	docker-compose restart api

docker-clean:
	docker-compose down -v  # removes volumes too — wipes DB data

psql:
	docker exec -it ecommerce_db psql -U postgres -d e-commerce

build:
	CGO_ENABLED=0 go build -o bin/server ./cmd/api/main.go

test:
	go test ./... -v
