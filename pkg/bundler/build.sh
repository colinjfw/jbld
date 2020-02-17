#!/bin/bash

runtime=$(./node_modules/.bin/uglifyjs pkg/bundler/runtime.js)
echo "package bundler" > pkg/bundler/runtime.go
echo "" >> pkg/bundler/runtime.go
echo "const runtime = \`${runtime}" >> pkg/bundler/runtime.go
echo "\`" >> pkg/bundler/runtime.go
