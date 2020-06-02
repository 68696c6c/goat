.DEFAULT:
	@echo 'App targets:'
	@echo
	@echo '    build   build api image and compile the app'
	@echo '    deps    install dependancies'
	@echo '    test    run unit tests'
	@echo


default: .DEFAULT

build:
	go build -i -o goat

deps:
	go mod tidy
	go mod vendor

test:
	go test ./... -cover
