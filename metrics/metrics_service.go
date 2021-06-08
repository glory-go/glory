package metrics

import "github.com/glory-go/glory/config"

type MetricsService interface {
	loadconfig(conf *config.MetricsConfig) error
	setup()
	run()
	counterInsc(metrics string)
	gaugeSet(gauge string, val float64)
	getMetrisServiceName() string
}

// MetricsServiceConfig 注释
// serverName == jobName == 界面上看到的job=""
// 如果当前service配置文件中有jobName则针对当前service使用当前jobName作为serverName，在前端统一打点服务中有用到
// 如果配置文件service没有jobName，则使用配置文件的serverName作为当前service的serverName
// orgName 组织名：GoOnline 或者 IDE 或者 Classroom 或者 Children
// serverName 比如： project-service
type MetricsServiceConfig struct {
	clientPort         string
	clientPath         string
	metricsType        string
	actionType         string
	gatewayPort        string
	gatewayHost        string
	metricsServiceName string
}
