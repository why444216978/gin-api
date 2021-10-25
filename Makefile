PWD     := $(shell pwd)
BINARY  := $(shell basename `pwd`)
GOROOT  := $(shell echo /usr/local/go)
GO      := $(GOROOT)/bin/go
GOMOD   := $(GO) mod
GOBUILD := $(GO) build
GOTEST  := $(GO) test

all: build test

build:
	$(GOMOD) download 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) --ldflags '-extldflags "-static"' -o $(BINARY)

test:
	$(GOTEST) -race -cover -coverpkg=./... ./...  -gcflags="-N -l"

clean:
	rm -rf $(APPNAME)

.PHONY: all build test clean
