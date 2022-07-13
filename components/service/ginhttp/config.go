package ginhttp

type ginHttpServiceConfig struct {
	Addr         string `mapstructure:"addr"`
	ReadTimeout  int    `mapstructure:"read_timeout"`  // 单位：s
	WriteTimeout int    `mapstructure:"write_timeout"` // 单位：s
}
