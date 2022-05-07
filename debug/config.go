package debug

type Config struct {
	Enable bool   `yaml:"enable"`
	Port   string `yaml:"port"`
}

func (b *Config) Prefix() string {
	return "boot"
}
