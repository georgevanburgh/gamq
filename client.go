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
