install:
	make -C lib install
	make -C test/react-app install
	make -C test/basic install
.PHONY: install

build:
	@make -C lib build &> /dev/null
	@echo "built"
.PHONY: build

test/test:
	@make -C test/basic test &> /dev/null
	@echo "ok      github.com/colinjfw/jbld/test/basic"
	@make -C test/react-app test &> /dev/null
	@echo "ok      github.com/colinjfw/jbld/test/react-app"
.PHONY: test/test

test/pkg:
	@go test -cover ./pkg/...
.PHONY: test/pkg

test: build test/pkg test/test
.PHONY: test
