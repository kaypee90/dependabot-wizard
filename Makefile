GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
BINARY_NAME=depbot
BINARY_UNIX=$(BINARY_NAME)_unix
UI_DIR=ui/depbot-app

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

.PHONY: run help version web build test test-go test-ui typecheck verify-assets ui clean build-run all deps build-linux

run:
	$(GORUN) $(LDFLAGS) ./cmd/depbot
help:
	$(GORUN) ./cmd/depbot --help
version:
	$(GORUN) $(LDFLAGS) ./cmd/depbot --version
web:
	$(GORUN) $(LDFLAGS) ./cmd/depbot --web
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/depbot

## test runs both suites; the Go one alone will not catch a frontend regression.
test: test-go test-ui verify-assets
test-go:
	$(GOTEST) -v ./...
test-ui:
	cd $(UI_DIR) && CI=true npx react-scripts test --watchAll=false
typecheck:
	cd $(UI_DIR) && npx tsc --noEmit
verify-assets:
	./scripts/verify-embedded-assets.sh

## ui rebuilds the React app and re-embeds it in the Go binary.
ui:
	cd $(UI_DIR) && ./build.sh

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
build-run: build
	./$(BINARY_NAME)
all: test build run
deps:
	$(GOGET) github.com/manifoldco/promptui
	cd $(UI_DIR) && npm install --legacy-peer-deps

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) ./cmd/depbot
