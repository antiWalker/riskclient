GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=riskclient
BINARY_UNIX=$(BINARY_NAME)_unix

all: build
build:
	go env -w GOPROXY=https://goproxy.io,direct
	go env -w GO111MODULE=on
	go build -v -o $(BINARY_NAME) main.go
build-linux:
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
