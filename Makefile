CMD = cmd/main.go

run:
	go run $(CMD)
build:
	go build $(CMD)
lint:
	 golangci-lint run --enable-all