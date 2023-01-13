package grpc

type grpcConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`

	TLS                bool   `mapstructure:"tls"`
	CertPath           string `mapstructure:"cert_path"`            // 绝对路径
	ServerNameOverride string `mapstructure:"server_name_override"` // 仅用于测试。如果设置为非空字符串，它将覆盖请求中授权的虚拟主机名（例如：authority 标头字段）
}
