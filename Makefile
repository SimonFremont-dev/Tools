.PHONY: setup test lint fmt fmt-check security build smoke clean

BIN ?= bin/tools
SMOKE_DIR ?= /tmp/tools-smoke
STAGES ?= dev,quality,build,cicd,kubernetes,ai,docs

setup:
	go mod download

test:
	go test ./...

lint:
	golangci-lint run ./...

fmt:
	golangci-lint fmt ./...

fmt-check:
	golangci-lint fmt --diff ./...

security:
	govulncheck ./...

build:
	go build -o $(BIN) ./cmd/tools

smoke: build
	rm -rf $(SMOKE_DIR)
	$(BIN) new demo --lang go --module github.com/me/demo --stages $(STAGES) --dir $(SMOKE_DIR)
	cd $(SMOKE_DIR) && go build ./... && go run ./cmd/demo

clean:
	rm -rf bin $(SMOKE_DIR)
