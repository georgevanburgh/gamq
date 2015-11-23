package gamq

import (
	log "github.com/cihub/seelog"
	"sync/atomic"
	"time"
)

type DummyMessageHandler struct {
	messagesSentLastSecond uint64 // messagesSentLastSecond should never be > 0
	running                bool
}

func (dmh *DummyMessageHandler) Initialize(input chan *string) chan *string {

	outputChannel := make(chan *string)
	dmh.running = true
	go dmh.logMetrics()

	go func() {
		for {
			message, more := <-input
			atomic.AddUint64(&dmh.messagesSentLastSecond, 1)
			if more {
				// log.Debugf("%s\n", &message)
				outputChannel <- message
			} else {
				// We're done
				log.Info("Closing output channel")
				dmh.running = false
				close(outputChannel)
				return
			}
		}
	}()

	return outputChannel
}

func (dmh *DummyMessageHandler) logMetrics() {
	for range time.Tick(time.Second) {
		// Die with the handler
		if !dmh.running {
			break
		}

		// Print out the number of messages published last second
		currentValue := atomic.LoadUint64(&dmh.messagesSentLastSecond)
		log.Warnf("%d/second\n", currentValue)
		atomic.StoreUint64(&dmh.messagesSentLastSecond, 0)
	}
}
