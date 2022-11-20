include .env
export $(shell sed 's/=.*//' .env)
.PHONY: run

run: docker-deps
	go run main.go --config-dir "$(CURDIR)/config" --migration-dir "$(CURDIR)/migrations"

build:
	go build -v ./...

test: docker-deps
	go test ./leagues ./managers ./players ./stadiums ./teams

docker-deps:
	docker-compose up -d --build
