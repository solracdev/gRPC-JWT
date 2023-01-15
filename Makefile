build-server:
	@go build -o bin/server ./cmd/server/main.go

build-client:
	@go build -o bin/client ./cmd/client/main.go

run-server:
	@./bin/server

run-client:
	@./bin/client

test:
	@go test -v ./...