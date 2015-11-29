package gamq

import (
	log "github.com/cihub/seelog"
	"github.com/quipo/statsd"
	"time"
)

type MetricsManager struct {
	metricsChannel chan *Metric
	statsBuffer    *statsd.StatsdBuffer
	statsdEnabled  bool
}

func (m *MetricsManager) Initialize() chan<- *Metric {
	m.metricsChannel = make(chan *Metric, 100)

	if Configuration.StatsDEndpoint != "" {
		m.statsdEnabled = true
		statsClient := statsd.NewStatsdClient(Configuration.StatsDEndpoint, "gamq.")
		statsClient.CreateSocket()
		interval := time.Second
		m.statsBuffer = statsd.NewStatsdBuffer(interval, statsClient)
		log.Debug("Initialized StatsD")
	} else {
		m.statsdEnabled = false
	}

	go m.listenForMetrics()

	return m.metricsChannel
}

func (m *MetricsManager) listenForMetrics() {
	defer m.statsBuffer.Close()
	var metric *Metric
	for {
		metric = <-m.metricsChannel
		log.Debugf("Received metric: %s - %d", metric.Name, metric.Value)

		if m.statsdEnabled {
			log.Debugf("Logging metrics")
			switch metric.Type {
			case "counter":
				m.statsBuffer.Incr(metric.Name, metric.Value)
			case "guage":
				m.statsBuffer.Gauge(metric.Name, metric.Value)
			}
		}
	}
}
