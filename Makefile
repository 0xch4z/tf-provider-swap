PROJECT := tf-provider-swap

DIST_DIR := ./dist
BIN      := $(DIST_DIR)/$(PROJECT)

ENTRYPOINT := cmd/$(PROJECT)/main.go

default: build

fmt:
	go fmt ./...

test: fmt
	go test -v ./...

build: test
	mkdir -p $(DIST_DIR)
	go build -o $(BIN) $(ENTRYPOINT)
