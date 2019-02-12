
.PHONY: all test build build-linux-amd64 build-osx-amd64 build-windows-amd64 build-win-386

all: test build build-linux-amd64 build-osx-amd64 build-windows-amd64 build-win-386

test:
	@go test -v -race ./...

build:
	@go build -o gtw

build-linux-amd64:
	@GOOS=linux GOARCH=amd64 go build -o gtw-linux-amd64

build-linux-386:
	@GOOS=linux GOARCH=386 go build -o gtw-linux-386

build-osx-amd64:
	@GOOS=darwin GOARCH=amd64 go build -o gtw-osx-amd64

build-osx-386:
	@GOOS=darwin GOARCH=386 go build -o gtw-osx-386

build-windows-amd64:
	@GOOS=darwin GOARCH=amd64 go build -o gtw-win-amd64.exe

build-win-386:
	@GOOS=darwin GOARCH=386 go build -o gtw-win-386.exe