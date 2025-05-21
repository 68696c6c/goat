DCR = docker compose run --rm

NETWORK_NAME ?= docker-dev
APP_NAME = web
DB_NAME = db

.PHONY: docs build

.DEFAULT:
	@echo 'App targets:'
	@echo
	@echo '    image-local        build the $(APP_NAME):dev Docker image for local development'
	@echo '    image-built        build the $(APP_NAME):built Docker image for task running'
	@echo '    build              compile the app for use in Docker'
	@echo '    deps               install dependencies'
	@echo '    setup-network      create local Docker network'
	@echo '    setup              set up local databases'
	@echo '    local              start local dev environment'
	@echo '    local-down         tear down local dev environment'
	@echo '    test               run unit tests'
	@echo '    migration          create a new migration'
	@echo '    migrate            migrate the local database'
	@echo '    migrate-reset      migrate the local database down to 0 and back up to the latest version'
	@echo '    seed               resets the local database and seeds it with example data'
	@echo '    lint               run the linter'
	@echo '    lint-fix           run the linter and fix any problems'
	@echo


default: .DEFAULT

image-local:
	docker build . --target dev -t $(APP_NAME):dev

image-built:
	docker build . --target dev -t $(APP_NAME):built

build:
	$(DCR) app go build -o $(APP_NAME)

deps:
	$(DCR) app sh -c "go mod tidy && go mod vendor"

setup-network:
	docker network create docker-dev || exit 0

setup: setup-network image-local deps build
	$(DCR) app ./$(APP_NAME) migrate up

local: local-down build
	NETWORK_NAME="$(NETWORK_NAME)" docker compose up app

local-down:
	NETWORK_NAME="$(NETWORK_NAME)" docker compose down

test:
	$(DCR) app wait-for-it -s -t 60 db:3306 -- go test ./... -cover -coverprofile=../cover.out

coverage:
	go tool cover -html=cover.out

migration: build
	$(DCR) app goose -dir db/migrations create $(name)

migrate: build
	$(DCR) app ./$(APP_NAME) migrate up

migrate-down: build
	$(DCR) app ./$(APP_NAME) migrate down-to 0

migrate-reset: build
	$(DCR) app ./$(APP_NAME) migrate down-to 0
	$(DCR) app ./$(APP_NAME) migrate up

seed: migrate-reset
	$(DCR) app ./$(APP_NAME) seed

lint:
	$(DCR) app golangci-lint run -v

lint-fix:
	$(DCR) app golangci-lint run --fix
