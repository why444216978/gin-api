.PHONY: all build test clean

all: build test

format:
	go vet ./...
	gofmt -w .
	golint ./...
build:
	go mod download 

test:
	go test -race -cover -coverpkg=./... ./...  -gcflags="-N -l"

clean:
	go clean -i -n -r

