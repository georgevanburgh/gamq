package gamq

import (
	log "github.com/cihub/seelog"
)

type DummyMessageHandler struct {
	running bool
}

func (dmh *DummyMessageHandler) Initialize(input <-chan *string) chan *string {

	outputChannel := make(chan *string)
	dmh.running = true

	go func() {
		for {
			message, more := <-input
			if more {
				// log.Debugf("%s", &message)
				outputChannel <- message
			} else {
				// We're done
				log.Debug("Closing output channel")
				dmh.running = false
				close(outputChannel)
				return
			}
		}
	}()

	return outputChannel
}
