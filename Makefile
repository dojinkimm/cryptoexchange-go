GOPATH:=$(shell go env GOPATH)
APP?=cryptoexchange-go

.PHONY: format
format:
	gofmt -s -w .
	go mod tidy

.PHONY: test
test:
	@go install github.com/rakyll/gotest@latest
	gotest -p 1 -race -cover -v ./...

.PHONY: lint
lint:
	@go install github.com/kyoh86/scopelint@latest
	golangci-lint run ./...
	scopelint --set-exit-status ./...
	go mod verify