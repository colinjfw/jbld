#!/bin/bash
set -e

mkdir -p ./bin && rm -r ./bin && mkdir -p ./bin

TARGETS="darwin-amd64 windows-amd64 linux-amd64"

for target in $TARGETS; do
  os=$(echo $target | cut -f1 -d-)
  arch=$(echo $target | cut -f2 -d-)
  echo "building jbld-$os-$arch"
  GOOS=$os GOARCH=$arch go build -ldflags="-s -w" -o ./bin/$os-$arch .
done
