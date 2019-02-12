
.PHONY: all test build build-linux-x86_64 build-linux-x86 build-osx-x86_64 build-osx-x86 build-win-x86_64 build-win-x86

all: test build build-linux-x86_64 build-linux-x86 build-osx-x86_64 build-osx-x86 build-win-x86_64 build-win-x86

test:
	@echo "Running tests"
	@go test -v -race ./...

build:
	@echo "Building"
	@go build -o gtw

build-linux-x86_64:
	@echo "Building linux-x86_64 version"
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o gtw-linux-x86_64

build-linux-x86:
	@echo "Building linux-x86 version"
	@GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o gtw-linux-x86

build-osx-x86_64:
	@echo "Building osx-x86_64 version"
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o gtw-osx-x86_64

build-osx-x86:
	@echo "Building osx-x86 version"
	@GOOS=darwin GOARCH=386 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o gtw-osx-x86

build-win-x86_64:
	@echo "Building win-x86_64 version"
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o gtw-win-x86_64.exe

build-win-x86:
	@echo "Building win-x86 version"
	@GOOS=darwin GOARCH=386 CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -o gtw-win-x86.exe