package mysql

import "github.com/glory-go/glory/v2/config"

func init() {
	config.RegisterComponent(GetMysqlGormComponent())
}
