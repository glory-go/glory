package grpc

type grpcServiceConfig struct {
	Addr string `mapstructure:"addr"`

	TLS      bool   `mapstructure:"tls"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}
