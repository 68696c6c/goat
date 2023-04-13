DCR = docker-compose run --rm

NETWORK_NAME ?= docker-dev
APP_NAME = cli
DB_NAME = db

.PHONY: docs build

.DEFAULT:
	@echo 'App targets:'
	@echo
	@echo '    image       build the $(APP_NAME):dev Docker image'
	@echo '    build       compile the app for use in Docker'
	@echo '    deps        install dependencies'
	@echo '    setup       set up local environment'
	@echo '    test        run unit tests'
	@echo '    lint        run the linter'
	@echo '    lint-fix    run the linter and fix any problems'
	@echo '    person      fetch and print a Star Wars API person'
	@echo '    planet      fetch and print a Star Wars API planet'
	@echo '    starship    fetch and print a Star Wars API starship'
	@echo


default: .DEFAULT

image:
	docker build . -t $(APP_NAME):dev

build:
	$(DCR) app go build -o $(APP_NAME)

deps:
	$(DCR) app sh -c "go mod tidy && go mod vendor"

setup: image deps build

test:
	$(DCR) app go test ./... -cover -coverprofile=cover.out

coverage:
	go tool cover -html=src/cover.out

lint:
	$(DCR) app golangci-lint run -v

lint-fix:
	$(DCR) app golangci-lint run --fix

person: build
	$(DCR) app ./cli fetch-person 1

planet: build
	$(DCR) app ./cli fetch-planet 3

starship: build
	$(DCR) app ./cli fetch-starship 9
