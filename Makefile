all: build
test: build
	go test ./...

int_test: test
	go test ./... -tags=integration

build:
	go build ./...
