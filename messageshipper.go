package gamq

import (
	log "github.com/cihub/seelog"
)

type MessageShipper struct {
	messageChannel chan *string
	subscriber     *Client
	ClientName     string
	CloseChannel   chan bool
}

func (shipper *MessageShipper) Initialize(inputChannel chan *string, subscriber *Client) {
	shipper.subscriber = subscriber
	shipper.messageChannel = inputChannel
	shipper.CloseChannel = make(chan bool)
	shipper.ClientName = subscriber.Name

	go shipper.forwardMessageToClient()
}

func (shipper *MessageShipper) forwardMessageToClient() {
	for {
		select {
		case message, more := <-shipper.messageChannel:
			if more {
				_, err := shipper.subscriber.Writer.WriteString(*message)
				if err != nil {
					log.Errorf("Error whilst sending message to consumer: %s", err)
				}
				shipper.subscriber.Writer.Flush()
			} else {
				return
			}
		case closing := <-shipper.CloseChannel:
			if closing {
				log.Debugf("Message shipper for %s closing down", shipper.ClientName)
				return
			}
		}
	}
}
