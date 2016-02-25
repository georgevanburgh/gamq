package gamq

import (
	"github.com/FireEater64/gamq/message"
	log "github.com/cihub/seelog"
)

type DummyMessageHandler struct {
	running bool
}

func (dmh *DummyMessageHandler) Initialize(input <-chan *message.Message) chan *message.Message {

	outputChannel := make(chan *message.Message, 10000) // TODO: Should definitely be configurable
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
