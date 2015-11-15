package gamq

import (
	log "github.com/cihub/seelog"
)

type MessageShipper struct {
	subscriberChannel chan *Client
	messageChannel    chan string
	subscribers       []*Client
}

func (shipper *MessageShipper) Initialize(inputChannel chan string, subscriberChannel chan *Client) chan<- string {
	shipper.subscribers = make([]*Client, 0)
	shipper.subscriberChannel = subscriberChannel
	shipper.messageChannel = inputChannel

	go shipper.listenForNewSubscribers()
	go shipper.forwardMessageToClients()

	return shipper.messageChannel
}

func (shipper *MessageShipper) listenForNewSubscribers() {
	for {
		newClient, more := <-shipper.subscriberChannel
		if more {
			log.Info("New subscriber!\n")
			_ = "breakpoint"
			shipper.subscribers = append(shipper.subscribers, newClient)
		} else {
			return
		}
	}
}

func (shipper *MessageShipper) forwardMessageToClients() {
	for {
		message, more := <-shipper.messageChannel
		if more {
			_ = "breakpoint"
			for _, subscriber := range shipper.subscribers {
				subscriber.Writer.WriteString(message)
				subscriber.Writer.Flush()
			}
		} else {
			return
		}

	}
}
