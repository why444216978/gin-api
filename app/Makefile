.PHONY: all build test clean

all: build test

format:
	go vet ./...
	gofmt -w .
	golint ./...
build:
	go mod download 
	GOOS=linux GOARCH=386 go build -o app main.go

test:
	go test -race -cover -coverpkg=./... ./...  -gcflags="-N -l"

clean:
	go clean -i -n -r

