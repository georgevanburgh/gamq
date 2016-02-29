package gamq

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"testing"

	"github.com/onsi/gomega"
)

func TestMetricsManager_ReceivesBasicMetric_PublishesDownstreamAndSendsToStatsD(t *testing.T) {
	gomega.RegisterTestingT(t)

	// Listen on UDP
	var statsDBuffer [2048]byte
	var udpPacketsReceived int
	udpAddr, _ := net.ResolveUDPAddr("udp", ":0")
	udpConn, _ := net.ListenUDP("udp", udpAddr)

	// Don't care about the contents of the received messages - just the fact
	// that we received them. We trust the StatsD library
	go func() {
		for i := 0; i < 3; i++ {
			_, _, _ = udpConn.ReadFromUDP(statsDBuffer[0:])
			udpPacketsReceived++
		}
	}()

	config := Config{StatsDEndpoint: fmt.Sprintf("localhost:%d", udpConn.LocalAddr().(*net.UDPAddr).Port)}
	SetConfig(&config)

	qm := newQueueManager()

	// Listen to metrics queue
	writerBuffer := new(bytes.Buffer)
	dummyWriter := bufio.NewWriter(writerBuffer)
	closedChannel := make(chan bool)
	dummyClient := Client{Name: "Test", Writer: dummyWriter, Closed: &closedChannel}

	qm.Subscribe("metrics", &dummyClient)

	// Log one of each metric
	// Check we've received metrics both via UDP - and on the metrics channel
	testMetric := NewMetric("test", "guage", 123)
	qm.metricsManager.metricsChannel <- testMetric

	gomega.Eventually(func() int {
		return udpPacketsReceived
	}, "2s").Should(gomega.Equal(1))

	gomega.Eventually(func() []byte {
		return writerBuffer.Bytes()
	}).ShouldNot(gomega.BeNil())

	writerBuffer.Reset()

	testMetric2 := NewMetric("test", "counter", 123)
	qm.metricsManager.metricsChannel <- testMetric2

	gomega.Eventually(func() int {
		return udpPacketsReceived
	}, "2s").Should(gomega.Equal(2))

	gomega.Eventually(func() []byte {
		return writerBuffer.Bytes()
	}).ShouldNot(gomega.BeNil())

	writerBuffer.Reset()

	testMetric3 := NewMetric("test", "timing", 123)
	qm.metricsManager.metricsChannel <- testMetric3

	gomega.Eventually(func() int {
		return udpPacketsReceived
	}, "2s").Should(gomega.Equal(3))

	gomega.Eventually(func() []byte {
		return writerBuffer.Bytes()
	}).ShouldNot(gomega.BeNil())

	writerBuffer.Reset()
}
