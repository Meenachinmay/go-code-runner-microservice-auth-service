.PHONY: proto
proto:
		protoc --go_out=. --go-grpc_out=. proto/company_auth/v1/company_auth.proto

.PHONY: build
build:
	go build -o bin/server cmd/server/main.go

.PHONY: run
run:
	go run cmd/server/main.go

.PHONY: test
test:
	go test -v ./...
