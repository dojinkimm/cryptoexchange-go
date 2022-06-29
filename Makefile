GOPATH:=$(shell go env GOPATH)
APP?=cryptox-go

.PHONY: format
## format: format files
format:
	gofmt -s -w .
	go mod tidy

.PHONY: test
## test: run tests
test:
	@go get github.com/rakyll/gotest
	gotest -p 1 -race -cover -v ./...

.PHONY: lint
## lint: check everything's okay
lint:
	@go get github.com/kyoh86/scopelint
	golangci-lint run ./...
	scopelint --set-exit-status ./...
	go mod verify