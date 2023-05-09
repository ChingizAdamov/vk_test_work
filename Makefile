.SILENT:

build:
	GOPATH=$(PWD) go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/main