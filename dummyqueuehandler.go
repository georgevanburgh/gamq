package gamq

import (
	log "github.com/cihub/seelog"
	"sync/atomic"
	"time"
)

type DummyMessageHandler struct {
	messagesSentLastSecond uint64
}

func (dmh *DummyMessageHandler) Initialize(input chan string) chan string {
	outputChannel := make(chan string)
	go dmh.logMetrics()
	go func() {
		for {
			message, more := <-input
			atomic.AddUint64(&dmh.messagesSentLastSecond, 1)
			if more {
				log.Debugf("%s\n", message)
				outputChannel <- message
			} else {
				// We're done
				log.Info("Closing output channel")
				close(outputChannel)
				return
			}
		}
	}()
	return outputChannel
}

func (dmh *DummyMessageHandler) logMetrics() {
	for range time.Tick(time.Second) {
		currentValue := atomic.LoadUint64(&dmh.messagesSentLastSecond)
		log.Infof("%d/second\n", currentValue)
		dmh.messagesSentLastSecond = 0
	}
}
