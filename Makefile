.DEFAULT:
	@echo 'App targets:'
	@echo
	@echo '    test    run unit tests'
	@echo


default: .DEFAULT

test:
	go test ./... -cover
