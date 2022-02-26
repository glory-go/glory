package mysql

type MysqlConfig struct {
	Host     string `yaml:"host" json:"host,omitempty"`
	Port     string `yaml:"port" json:"port,omitempty"`
	Username string `yaml:"username" json:"username,omitempty"`
	Password string `yaml:"password" json:"password,omitempty"`
	DBName   string `yaml:"dbname" json:"dbname,omitempty"`
}
