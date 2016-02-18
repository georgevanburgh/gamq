package gamq

import (
	"bufio"
)

type Client struct {
	Name   string
	Writer *bufio.Writer
	Reader *bufio.Reader
	Closed *chan bool
}

func NewClient(givenName string, givenWriter *bufio.Writer, givenReader *bufio.Reader) *Client {
	closedChannel := make(chan bool)
	return &Client{
		Name:   givenName,
		Writer: givenWriter,
		Reader: givenReader,
		Closed: &closedChannel}
}
