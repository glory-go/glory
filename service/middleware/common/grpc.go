package common

import "google.golang.org/grpc"

var (
	unaryClientMWs []grpc.UnaryClientInterceptor
	unaryServerMWs []grpc.UnaryServerInterceptor
)

func RegisterUnaryClientMWs(filters ...grpc.UnaryClientInterceptor) {
	unaryClientMWs = append(unaryClientMWs, filters...)
}

func RegisterUnaryServerMWs(filters ...grpc.UnaryServerInterceptor) {
	unaryServerMWs = append(unaryServerMWs, filters...)
}

func GetUnaryClientMWs() []grpc.UnaryClientInterceptor {
	return unaryClientMWs
}

func GetUnaryServerMWs() []grpc.UnaryServerInterceptor {
	return unaryServerMWs
}
