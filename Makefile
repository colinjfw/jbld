install:
	make -C lib install
	make -C examples/react-app install
	make -C examples/basic install
.PHONY: install

build:
	@make -C lib build &> /dev/null
.PHONY: build

test/examples:
	@make -C examples/basic test &> /dev/null
	@echo "ok      github.com/colinjfw/jbld/examples/basic"
	@make -C examples/react-app test &> /dev/null
	@echo "ok      github.com/colinjfw/jbld/examples/react-app"
.PHONY: test/examples

test/pkg:
	@go test -cover ./pkg/...
.PHONY: test/pkg

test: build test/pkg test/examples
.PHONY: test
