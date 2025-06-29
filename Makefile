# Makefile

APP_NAME=bondhupay
BINARY=tmp/$(APP_NAME)

.PHONY: run build clean test fmt

dev:
	air

build:
	go build -o $(BINARY) ./cmd/main.go

clean:
	rm -rf tmp/

test:
	go test ./...

fmt:
	go fmt ./...
