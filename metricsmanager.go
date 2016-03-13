package gamq

import (
	"fmt"
	"github.com/FireEater64/gamq/message"
	log "github.com/cihub/seelog"
	"github.com/quipo/statsd"
	"time"
)

const (
	metricsQueueName = "metrics"
)

type MetricsManager struct {
	metricsChannel chan *Metric
	statsBuffer    *statsd.StatsdBuffer
	statsdEnabled  bool
	queueManager   *queueManager
}

func NewMetricsManager(givenQueueManager *queueManager) *MetricsManager {
	m := MetricsManager{}

	m.queueManager = givenQueueManager

	m.metricsChannel = make(chan *Metric, 100)

	if Configuration.StatsDEndpoint != "" {
		m.statsdEnabled = true
		statsClient := statsd.NewStatsdClient(Configuration.StatsDEndpoint, "gamq.")
		statsClient.CreateSocket()
		interval := time.Second
		m.statsBuffer = statsd.NewStatsdBuffer(interval, statsClient)
		log.Debugf("Initialized StatsD - sending metrics to: %s", Configuration.StatsDEndpoint)
	} else {
		m.statsdEnabled = false
	}

	go m.listenForMetrics()

	return &m
}

func (m *MetricsManager) listenForMetrics() {
	if m.statsdEnabled {
		defer m.statsBuffer.Close()
	}

	var metric *Metric
	for {
		metric = <-m.metricsChannel
		log.Debugf("Received metric: %s - %v", metric.Name, metric.Value)

		if m.statsdEnabled {
			log.Debugf("Logging metrics")
			switch metric.Type {
			case "counter":
				m.statsBuffer.Incr(metric.Name, metric.Value)
			case "guage":
				m.statsBuffer.Gauge(metric.Name, metric.Value)
			case "timing":
				m.statsBuffer.Timing(metric.Name, metric.Value)
			default:
				log.Errorf("Unknown metric type received: %s", metric.Type)
			}
		}

		stringToPublish := fmt.Sprintf("%s:%s", metric.Name, metric.Value)
		messageHeaders := make(map[string]string)
		messageBody := []byte(stringToPublish)

		metricMessage := message.NewMessage(&messageHeaders, &messageBody)
		m.queueManager.Publish(metricsQueueName, metricMessage)
	}
}
