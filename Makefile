include .env
export $(shell sed 's/=.*//' .env)
.PHONY: run

run: docker-deps
	go run main.go --config-dir "$(CURDIR)/config" --migration-dir "$(CURDIR)/migrations"


docker-deps:
	docker-compose up -d --build
