package metrics

import (
	"net/http"

	gxnet "github.com/dubbogo/gost/net"
	"github.com/glory-go/glory/tools"

	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PromePullService struct {
	config         *MetricsServiceConfig
	metricsCounter *prometheus.CounterVec
	metricsGauge   *prometheus.GaugeVec
}

func (p *PromePullService) loadconfig(conf *config.MetricsConfig) error {
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

func (p *PromePullService) run() {
	log.Info("port = ", p.config.clientPort, "path = ", p.config.clientPath)
	http.Handle(p.config.clientPath, promhttp.Handler())
	log.Info(http.ListenAndServe(":"+p.config.clientPort, nil))
	log.Info("prometheus stop listening")
}

func (p *PromePullService) setup() {
	// prometheus pull模式依赖counter
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

	prometheus.MustRegister(p.metricsCounter)
	prometheus.MustRegister(p.metricsGauge)
}

func (p *PromePullService) counterInsc(metricsType string) {
	p.metricsCounter.With(prometheus.Labels{"metricsName": metricsType}).Inc()
	log.Debug("pull: counter ", metricsType, " insc")
}

func (p *PromePullService) gaugeSet(gauge string, val float64) {
	p.metricsGauge.With(prometheus.Labels{"metricsName": gauge}).Set(val)
	log.Debug("pull: gauge ", gauge, " set ", val)
}

func (p *PromePullService) getMetrisServiceName() string {
	return p.config.metricsServiceName
}

func NewPromePullService() *PromePullService {
	return &PromePullService{}
}
