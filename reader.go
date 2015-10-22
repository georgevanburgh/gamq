package gamq

import (
	"bufio"
	"fmt"
	"net"
)

const (
	PORT = 4321
)

type Reader interface {
	Initialize()
}

type TcpReader struct {
}

func (reader TcpReader) Initialize() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))

	if err != nil {
		fmt.Printf("An error occured whilst opening a socket for reading: %s",
			err.Error())
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("An error occured whilst opening a socket for reading: %s",
				err.Error())
		}

		channel := make(chan string)
		go listenForMessages(channel)
		go handleConnection(channel, &conn)
	}
}

func listenForMessages(in chan string) {
	for {
		stringToPrint := <-in
		fmt.Print(stringToPrint)
	}
}

func handleConnection(out chan string, conn *net.Conn) {

	bufferedReader := bufio.NewReader(*conn)

	for {
		line, err := bufferedReader.ReadBytes('\n')

		if err != nil {
			// Connection has been closed
			return
		}
		stringLine := string(line[:len(line)])
		out <- stringLine
	}
}
