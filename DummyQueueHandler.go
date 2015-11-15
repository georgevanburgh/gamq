package gamq

import (
	log "github.com/cihub/seelog"
)

type DummyMessageHandler struct {
}

func (dmh *DummyMessageHandler) Initialize(input chan string) chan string {
	outputChannel := make(chan string)
	go func() {
		for {
			message, more := <-input
			if more {
				log.Debugf("%s", message)
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
