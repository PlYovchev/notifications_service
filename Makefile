# Read the .env file and export its variables
ifeq (,$(wildcard ./.env))
$(error .env file not found)
else
include .env
export $(shell sed 's/=.*//' .env)
endif

PROJECT_NAME = $(shell basename "$(PWD)" | tr '[:upper:]' '[:lower:]')

# GIT commit id will be used as version of the application
VERSION ?= $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags "-X main.version=${VERSION}"

MODULE = $(shell go list -m)

## start: Starts everything that is required to serve the APIs
start:
	docker-compose up -d

## setup: Start the dependencies only
setup:
	docker-compose build

## run: Run the API server alone (without supplementary services such as DB etc.,)
run:
	go run ${LDFLAGS} main.go -version="${VERSION}"

## build: Build the API server binary
build:
	CGO_ENABLED=0 go build ${LDFLAGS} -a -o ${PROJECT_NAME} $(MODULE)

## version: Display the current version of the API server
version:
	@echo $(VERSION)

## test: Run tests
test:
	go test ./... -v -coverprofile coverage.out -covermode count

## tidy: Tidy go modules
tidy:
	go mod tidy

## format: Format go code
format:
	go fmt ./...

## lint: Run linter
lint:
	golangci-lint run

## lint-fix: Run linter and fix the issues
lint-fix:
	golangci-lint run --fix

## clean: Clean all docker resources
clean:
	docker-compose down
	docker ps --filter name=orders -q | xargs docker stop
	docker network ls --filter name=orders -q | xargs docker prune --force
	docker volume ls --filter name=orders -q | xargs docker volume rm

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command to run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |	sed -e 's/^/ /'
	@echo