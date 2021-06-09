package metrics

import (
	gxnet "github.com/dubbogo/gost/net"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/tools"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type PromePushService struct {
	config         *MetricsServiceConfig
	metricsCounter *prometheus.CounterVec
	metricsGauge   *prometheus.GaugeVec
	pusher         *push.Pusher
}

func (p *PromePushService) loadconfig(conf *config.MetricsConfig) error {
	var serverName string
	if conf.JobName != "" {
		serverName = conf.JobName
	} else {
		serverName = config.GlobalServerConf.ServerName
	}
	p.config = &MetricsServiceConfig{
		clientPort:         conf.ClientPort,
		metricsType:        conf.MetricsType,
		clientPath:         conf.ClientPath,
		actionType:         conf.ActionType,
		gatewayHost:        conf.GateWayHost,
		gatewayPort:        conf.GateWayPort,
		metricsServiceName: tools.PrometheusParseToSupportMetricsName(config.GlobalServerConf.OrgName + "_" + serverName),
	}
	return nil
}

func (p *PromePushService) run() {

}

func (p *PromePushService) setup() {
	localIP, _ := gxnet.GetLocalIP()
	p.metricsCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:        p.config.metricsServiceName + "_counter",
		ConstLabels: prometheus.Labels{"instance": localIP},
	}, []string{"metricsName"},
	)

	// prometheus pull 模式依赖gauge
	p.metricsGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:        p.config.metricsServiceName + "_gauge",
		ConstLabels: prometheus.Labels{"instance": localIP},
	}, []string{"metricsName"},
	)

	registry := prometheus.NewRegistry()
	registry.MustRegister(p.metricsCounter)
	registry.MustRegister(p.metricsGauge)
	p.pusher = push.New("http://"+p.config.gatewayHost+":"+p.config.gatewayPort, p.config.metricsServiceName).Gatherer(registry)
}

func (p *PromePushService) counterInsc(metricsType string) {
	p.metricsCounter.With(prometheus.Labels{"metricsName": metricsType}).Inc()
	if err := p.pusher.Push(); err != nil {
		log.Error("Could not push metrics to Pushgateway", err)
		return
	}
	log.Debug("push: counter ", metricsType, "insc")
}

func (p *PromePushService) gaugeSet(gauge string, val float64) {
	p.metricsGauge.With(prometheus.Labels{"metricsName": gauge}).Set(val)
	if err := p.pusher.Push(); err != nil {
		log.Error("Could not push metrics to Pushgateway", err)
		return
	}
	log.Debug("push: gauge ", gauge, " set ", val)
}

func (p *PromePushService) getMetrisServiceName() string {
	return p.config.metricsServiceName
}

func NewPromePushService() *PromePushService {
	return &PromePushService{}
}
