package metrics

import (
	"sync"

	"github.com/glory-go/glory/tools"

	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
)

type MetricsHandler interface {
	setup(config []*config.MetricsConfig)
	run()
	counterInsc(metricsName string, jobNames []string)
	gaugeSet(gaugeName string, val float64, jobNames []string)
}

type DefaultMetricsHandler struct {
	services []MetricsService
}

func (ph *DefaultMetricsHandler) setup(config []*config.MetricsConfig) {
	for _, v := range config {
		switch v.MetricsType {
		case "prometheus_client":
			var service MetricsService
			if v.ActionType == "pull" {
				service = NewPromePullService()
			} else if v.ActionType == "push" {
				service = NewPromePushService()
			} else {
				break
			}
			if err := service.loadconfig(v); err != nil {
				log.Error("serivice loadconfig err")
				return
			}
			service.setup()
			ph.services = append(ph.services, service)
		}
	}
}

func (ph *DefaultMetricsHandler) run() {
	wg := sync.WaitGroup{}
	for _, v := range ph.services {
		wg.Add(1)
		go func() {
			defer func() {
				if e := recover(); e != nil {
					log.Error("error %d")
					return
				}
				wg.Done()
			}()
			v.run()
		}()
	}
	log.Debug("prometheus service started")
	wg.Wait()
}

func (ph *DefaultMetricsHandler) counterInsc(metricsName string, jobNames []string) {
	if len(jobNames) == 0 {
		jobNames = append(jobNames, tools.PrometheusParseToSupportMetricsName(config.GlobalServerConf.OrgName+"_"+config.GlobalServerConf.ServerName))
	}

	haveSet := false
	jobNamesMap := make(map[string]bool)
	// 拿到当前所有需要发送的jobname
	for _, v := range jobNames {
		jobNamesMap[tools.PrometheusParseToSupportMetricsName(v)] = true
	}
	for _, v := range ph.services {
		// 如果需要发送到特定job
		if _, ok := jobNamesMap[v.getMetrisServiceName()]; ok {
			haveSet = true
			v.counterInsc(metricsName) // 通过当前job发送
		}
	}
	if !haveSet {
		log.Error("your jobNames = ", jobNames, "can't find target jobNames Service in glory.yaml, please check")
	}
	return
}

func (ph *DefaultMetricsHandler) gaugeSet(gaugeName string, val float64, jobNames []string) {
	if len(jobNames) != 0 {
		haveSet := false
		jobNamesMap := make(map[string]bool)
		for _, v := range jobNames {
			jobNamesMap[v] = true
		}
		for _, v := range ph.services {
			if _, ok := jobNamesMap[v.getMetrisServiceName()]; ok {
				haveSet = true
				v.gaugeSet(gaugeName, val)
			}
		}
		if !haveSet {
			log.Error("your jobNames = ", jobNames, "can't find target jobNames Service in glory.yaml, please check")
		}
		return
	}
	for _, v := range ph.services {
		v.gaugeSet(gaugeName, val)
	}
}

func NewDefaultMetricsHandler() MetricsHandler {
	return &DefaultMetricsHandler{}
}
