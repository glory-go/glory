proto-gen:
	protoc --go_out=./boot/api --go-grpc_out=./boot/api ./boot/api/glory/boot/debug.proto