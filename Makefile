GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
BINARY_NAME=depbot
BINARY_UNIX=$(BINARY_NAME)_unix

LAUNCH=$(GORUN) $$(ls -1 cmd/depbot/*.go | grep -v _test.go)

run:
	$(LAUNCH)
help:
	$(LAUNCH) --help
version:
	$(LAUNCH) --version
web:
	$(LAUNCH) --web
build: 
	$(GOBUILD) -o $(BINARY_NAME) $$(ls -1 cmd/depbot/*.go | grep -v _test.go)
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
build-run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
all: test build run
deps:
	$(GOGET) github.com/manifoldco/promptui

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
