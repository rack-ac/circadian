APP_NAME := circadian
MAIN_PKG := ./cmd/circadian
DIST_DIR := dist
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/rack-ac/circadian/internal/version.Version=$(VERSION)

.PHONY: build test clean dist

build:
	go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) $(MAIN_PKG)

test:
	go test ./...

dist:
	mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-darwin-arm64 $(MAIN_PKG)
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-darwin-amd64 $(MAIN_PKG)

clean:
	rm -rf $(APP_NAME) $(DIST_DIR)
