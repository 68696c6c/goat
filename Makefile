IMAGE_NAME = goat
BUILD_TAG ?= latest
APP_PATH = /go/goat

.PHONY: image dep local local-down test migrate

.DEFAULT:
	@echo 'Invalid target.'
	@echo
	@echo '    image           			build app image'
	@echo '    build           			build app image and compile the app'
	@echo '    dep            			install dependancies'
	@echo '    local          			spin up local environment'
	@echo '    local-down     			tear down local environment'
	@echo '    test           			run unit tests'
	@echo

default: .DEFAULT

image:
	docker build . -f Dockerfile -t $(IMAGE_NAME):$(BUILD_TAG)

build: image
	docker-compose run --rm app go build -i -o app

dep:
	docker-compose run --rm app dep ensure

local:
	docker-compose up

local-down:
	docker-compose down

test:
	docker-compose run --rm app go test . -cover
