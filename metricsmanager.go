package gamq

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/quipo/statsd"
	"time"
)

const (
	MetricsQueueName = "metrics"
)

type MetricsManager struct {
	metricsChannel chan *Metric
	statsBuffer    *statsd.StatsdBuffer
	statsdEnabled  bool
	queueManager   *QueueManager
}

func (m *MetricsManager) Initialize(givenQueueManager *QueueManager) chan<- *Metric {
	m.queueManager = givenQueueManager

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
	if m.statsdEnabled {
		defer m.statsBuffer.Close()
	}

	var metric *Metric
	for {
		metric = <-m.metricsChannel
		log.Debugf("Received metric: %s - %v", metric.Name, metric.Value)

		stringToPublish := fmt.Sprintf("%s:%s", metric.Name, metric.Value)
		m.queueManager.Publish(MetricsQueueName, &stringToPublish)

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
