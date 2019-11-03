IMAGE_NAME = goat
BUILD_TAG ?= latest
APP_PATH = /go/goat
CMD ?= default
TEST_APP_PATH ?= ~/Code/Go/src/github.com/68696c6c/goat-test

.PHONY: image dep cli local-down test migrate

.DEFAULT:
	@echo 'Invalid target.'
	@echo
	@echo '    image         build Docker image'
	@echo '    deps          install dependancies'
	@echo '    build         build the CLI for the current machine'
	@echo '    test          run unit tests'
	@echo '    new           generate a new Goat project'
	@echo

default: .DEFAULT

image:
	docker build . -f docker/Dockerfile -t $(IMAGE_NAME):$(BUILD_TAG)

deps:
	docker-compose run --rm app go mod tidy
	docker-compose run --rm app go mod vendor

build:
	 go build -o /usr/local/bin/goat-cli

cli-down:
	docker-compose down

test:
	docker-compose run --rm app go test ./... -cover

new: build
	rm -rf $(TEST_APP_PATH)
	goat-cli new goat-test.yml
