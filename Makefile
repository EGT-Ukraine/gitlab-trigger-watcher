
.PHONY: all test build

all: test build

test:
	@go test -v -race ./...

build:
	@go build