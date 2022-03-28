package mq

import "github.com/glory-go/glory/config"

var (
	inited bool
)

func Init() {
	loadConfig(config.GlobalServerConf.MQConfig)
	inited = true
}
