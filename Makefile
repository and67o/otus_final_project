BIN_ROTATION = "./bin"

build:
	go build -v -o ${BIN_ROTATION} ./cmd
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.33.0

lint: install-lint-deps
	golangci-lint run ./...
test:
	go test -race ./internal/...
generate:
	go generate ./internal/...
test-integration:
	go test -count=1 ./tests/...
