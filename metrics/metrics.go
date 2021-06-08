package metrics

import (
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
)

func init() {
	defaultMetricsHandler = NewDefaultMetricsHandler()
	defaultMetricsHandler.setup(config.GlobalServerConf.MetricsConfigs)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.Error("defaulteMetricsHandler run err = ", e)
			}
		}()
		defaultMetricsHandler.run()
	}()
}

// 默认metrics handler
var defaultMetricsHandler MetricsHandler

// 递增数据自增
func CounterInsc(metricsName string, jobNames ...string) {
	defaultMetricsHandler.counterInsc(metricsName, jobNames)
}

// 变动浮点数
func GaugeSet(gaugeName string, val float64, jobNames ...string) {
	defaultMetricsHandler.gaugeSet(gaugeName, val, jobNames)
}
