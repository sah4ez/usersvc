SOURCE=$(shell find ./* -type f -name '*.go' | xargs)
BUILD_PATH=./cmd/usersvc
GO=CC=gcc vgo

build: clean $(SOURCE)
	$(GO) build -o ./bin/usersvc $(BUILD_PATH)

vendor: go.mod
	$(GO) vendor

.PHONY: clean
clean: 
	rm -rf ./bin/
