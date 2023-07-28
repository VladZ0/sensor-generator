include .env.postgres

swagger:
	swag init -g ./cmd/main/main.go -o ./docs

test:
	go test ./internal/...

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

migrate:
	psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) --password $(POSTGRES_PASSWORD) -f ./migrations/1_init.up.sql

start_redis:
	docker run --name redis -d redis

stop_redis:
	docker stop redis
	docker rm redis