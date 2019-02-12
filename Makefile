
.PHONY: all test build build-linux-amd64 build-osx-amd64 build-windows-amd64 build-win-386

all: test build build-linux-amd64 build-osx-amd64 build-windows-amd64 build-win-386

test:
	@echo "Running tests"
	@go test -v -race ./...

build:
	@echo "Building"
	@go build -o gtw

build-linux-amd64:
	@echo "Building linux-amd64 version"
	@GOOS=linux GOARCH=amd64 go build -o gtw-linux-amd64

build-linux-386:
	@echo "Building linux-386 version"
	@GOOS=linux GOARCH=386 go build -o gtw-linux-386

build-osx-amd64:
	@echo "Building osx-amd64 version"
	@GOOS=darwin GOARCH=amd64 go build -o gtw-osx-amd64

build-osx-386:
	@echo "Building osx-386 version"
	@GOOS=darwin GOARCH=386 go build -o gtw-osx-386

build-windows-amd64:
	@echo "Building windows-amd64 version"
	@GOOS=darwin GOARCH=amd64 go build -o gtw-win-amd64.exe

build-win-386:
	@echo "Building win-386 version"
	@GOOS=darwin GOARCH=386 go build -o gtw-win-386.exe