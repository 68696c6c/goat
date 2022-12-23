.DEFAULT:
	@echo 'App targets:'
	@echo
	@echo '    deps    install dependencies'
	@echo '    test    run unit tests'
	@echo

default: .DEFAULT

deps:
	go mod tidy
	go mod vendor

test:
	go test ./... -cover
