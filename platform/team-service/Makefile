include .env
export $(shell sed 's/=.*//' .env)

.PHONY: build test run format

build:
	go build

test: docker-deps
	go test ./

run: docker-deps
	go run main.go

format:
	gofmt -s -w ./
	goimports -w -local github.com/upside-tech ./

docker-deps:
	$(MAKE) -C ../ docker-run
