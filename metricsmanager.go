package gamq

import (
	"fmt"
	"github.com/FireEater64/gamq/message"
	log "github.com/cihub/seelog"
)

const (
	metricsQueueName = "metrics"
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
		messageHeaders := make(map[string]string)
		messageBody := []byte(stringToPublish)

		metricMessage := message.NewMessage(&messageHeaders, &messageBody)
		m.queueManager.Publish(metricsQueueName, metricMessage)
	}
}
