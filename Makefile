GOPATH:=$(shell go env GOPATH)
APP?=cryptoexchange-go

.PHONY: format
format:
	gofmt -s -w .
	go mod tidy

.PHONY: test
test:
	@go get github.com/rakyll/gotest
	gotest -p 1 -race -cover -v ./...

.PHONY: lint
lint:
	@go install github.com/kyoh86/scopelint@latest
	golangci-lint run ./...
	scopelint --set-exit-status ./...
	go mod verify