proto-gen:
	protoc --go_out=./debug/api --go-grpc_out=./debug/api ./debug/api/glory/boot/debug.proto

tidy:
	go mod tidy

test:
	go test ./...

imports:
	imports-formatter .

lint: tidy
	golangci-lint run