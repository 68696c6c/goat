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

#test:
#	go test ./... -cover

image:
	docker build . --target dev -t goat:test

build:
	go build

test:
	docker-compose run --rm test

coverage:
	go tool cover -html=cover.out
