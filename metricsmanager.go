package gamq

import (
	"fmt"
	log "github.com/cihub/seelog"
)

const (
	MetricsQueueName = "metrics"
)

type MetricsManager struct {
	metricsChannel chan *Metric
	queueManager   *QueueManager
}

func (m *MetricsManager) Initialize(givenQueueManager *QueueManager) chan<- *Metric {
	m.queueManager = givenQueueManager

	m.metricsChannel = make(chan *Metric, 100)

	go m.listenForMetrics()

	return m.metricsChannel
}

func (m *MetricsManager) listenForMetrics() {
	var metric *Metric
	for {
		metric = <-m.metricsChannel
		log.Debugf("Received metric: %s - %v", metric.Name, metric.Value)

		stringToPublish := fmt.Sprintf("%s:%s", metric.Name, metric.Value)
		m.queueManager.Publish(MetricsQueueName, &stringToPublish)
	}
}
