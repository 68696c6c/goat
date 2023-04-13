.DEFAULT:
	@echo 'App targets:'
	@echo
	@echo '    deps    install dependencies'
	@echo '    test    run unit tests'
	@echo '    build   build the package'
	@echo

default: .DEFAULT

deps:
	make -C src deps

test:
	make -C src test

build:
	make -C src build
