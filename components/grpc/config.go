package grpc

type grpcConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
