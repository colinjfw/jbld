install:
	yarn install
.PHONY: install

test:
	../../lib/cli.js --clean
	../../lib/cli.js
.PHONY: test

start:
	make -C ../../lib build
	../../lib/cli.js --watch --serve ":3000"
.PHONY: start

build:
	../../lib/cli.js --mode production
.PHONY: build
