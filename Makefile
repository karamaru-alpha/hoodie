include .env
GO_VERSION := $(shell grep -m 1 -o '[0-9]\+\.[0-9]\+\.[0-9]\+' go.mod)
export

.PHONY: install
install:
	go mod download
	go install golang.org/x/tools/cmd/goimports@v0.21.0
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.0

.PHONY: buf-gen
buf-gen:
	docker compose run --rm buf generate ./proto/rpc --template ./proto/buf.gen.yaml
	make fmt

.PHONY: fmt
fmt:
	goimports -w -local "github.com/karamaru-alpha/days" pkg/
	gofmt -s -w pkg/

.PHONY: run-api
run-api:
	go run cmd/api/main.go

.PHONY: lint
lint:
	golangci-lint run -v cmd/... pkg/... config/...
