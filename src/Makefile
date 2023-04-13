.DEFAULT:
	@echo 'App targets:'
	@echo
	@echo '    deps    install dependencies'
	@echo '    test    run unit tests'
	@echo '    build   build the package'
	@echo

default: .DEFAULT

deps:
	go mod tidy
	go mod vendor

build:
	go build

test:
	docker-compose run --rm test

coverage:
	go tool cover -html=cover.out

image:
	docker build . --target dev -t goat:test

db:
	docker-compose run --rm db

down:
	docker-compose down
