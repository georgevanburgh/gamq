package gamq

import (
	"github.com/FireEater64/gamq/message"
	log "github.com/cihub/seelog"
)

type MessageShipper struct {
	messageChannel chan *message.Message
	subscriber     *Client
	ClientName     string
	CloseChannel   chan bool
}

func (shipper *MessageShipper) Initialize(inputChannel chan *message.Message, subscriber *Client) {
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
				_, err := shipper.subscriber.Writer.WriteString(string(*message.Body))
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
