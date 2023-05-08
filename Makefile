.SILENT:

build:
	GOPATH=$(PWD) go build -o ./.bin/main ./cmd/main

run: build
	./.bin/main