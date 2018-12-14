all: ready_to_commit

test: build
	go test ./...

int_test: test
	go test ./... -tags=integration

build:
	go build ./...

vet: build test 
	go vet ./...

ready_to_commit: build test int_test vet
