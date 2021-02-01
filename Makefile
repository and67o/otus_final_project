BIN_ROTATION = "./bin"

build:
	go build -v -o ${BIN_ROTATION} ./cmd
lint:
	golangci-lint run ./...