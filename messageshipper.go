package gamq

import (
	"github.com/FireEater64/gamq/message"
	log "github.com/cihub/seelog"
	"time"
)

type messageShipper struct {
	messageChannel chan *message.Message
	metricsChannel chan<- *Metric
	subscriber     *Client
	ClientName     string
	CloseChannel   chan bool
}

func newMessageShipper(inputChannel chan *message.Message, subscriber *Client, givenMetricsChannel chan<- *Metric) *messageShipper {
	shipper := messageShipper{}

	shipper.subscriber = subscriber
	shipper.messageChannel = inputChannel
	shipper.CloseChannel = make(chan bool)
	shipper.ClientName = subscriber.Name
	shipper.metricsChannel = givenMetricsChannel

	go shipper.forwardMessageToClient()

	return &shipper
}

func (shipper *messageShipper) forwardMessageToClient() {
	for {
		select {
		case message, more := <-shipper.messageChannel:
			if more {
				_, err := shipper.subscriber.Writer.WriteString(string(*message.Body))
				if err != nil {
					log.Errorf("Error whilst sending message to consumer: %s", err)
				}
				shipper.subscriber.Writer.Flush()

				// Calculate and log the latency for the sent message
				shipper.metricsChannel <- NewMetric("latency", "timing", time.Now().Sub(message.ReceivedAt).Nanoseconds()/1000000)
				// Log the number of bytes received
				shipper.metricsChannel <- NewMetric("bytesout.tcp", "count", int64(len(*message.Body)))

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
