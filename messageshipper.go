package gamq

import (
	log "github.com/cihub/seelog"
)

type MessageShipper struct {
	messageChannel chan *string
	subscribers    *[]*Client
}

func (shipper *MessageShipper) Initialize(inputChannel chan *string, subscriberArray *[]*Client) chan<- *string {
	shipper.subscribers = subscriberArray
	shipper.messageChannel = inputChannel

	go shipper.forwardMessageToClients()

	return shipper.messageChannel
}

func (shipper *MessageShipper) forwardMessageToClients() {
	for {
		message, more := <-shipper.messageChannel
		if more {
			_ = "breakpoint"
			for _, subscriber := range *shipper.subscribers {
				_, err := subscriber.Writer.WriteString(*message)
				if err != nil {
					log.Errorf("Error whilst sending message to consumer: %s", err)
				}
				subscriber.Writer.Flush()
			}
		} else {
			return
		}

	}
}
