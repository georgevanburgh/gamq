package gamq

import (
	"bufio"
	"fmt"
	log "github.com/cihub/seelog"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	TCPPORT                 = 48879
	HELPSTRING              = "Help. Me."
	UNRECOGNISEDCOMMANDTEXT = "Unrecognised command"
)

type ConnectionManager struct {
	wg   sync.WaitGroup
	qm   QueueManager
	rand *rand.Rand
}

func (manager *ConnectionManager) Initialize() {
	// Initialize our random number generator (used for naming new connections)
	// A different seed will be used on each startup, for no good reason
	manager.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	manager.qm = QueueManager{}
	manager.qm.Initialize()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", TCPPORT))

	if err != nil {
		log.Errorf("An error occured whilst opening a socket for reading: %s",
			err.Error())
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Errorf("An error occured whilst opening a socket for reading: %s",
				err.Error())
		}
		manager.wg.Add(1)
		go manager.handleConnection(&conn)
	}
}

func (manager *ConnectionManager) handleConnection(conn *net.Conn) {
	defer manager.wg.Done()
	connReader := bufio.NewReader(*conn)
	connWriter := bufio.NewWriter(*conn)
	client := Client{Name: strconv.Itoa(manager.rand.Int()),
		Writer: connWriter,
		Reader: connReader}

	for {
		// Read until newline
		line, err := client.Reader.ReadString('\n')

		if err != nil {
			// Connection has been closed
			break
		}

		stringLine := string(line[:len(line)])
		fmt.Print(stringLine)

		// Parse command and (optionally) return response (if any)
		manager.parseClientCommand(stringLine, &client)
	}

	log.Info("A connection was closed")
}

func (manager *ConnectionManager) parseClientCommand(command string, client *Client) {
	commandTokens := strings.Fields(command)

	if len(commandTokens) == 0 {
		return
	}

	switch strings.ToUpper(commandTokens[0]) {
	case "HELP":
		manager.sendStringToClient(HELPSTRING, client)
	case "PUB":
		manager.qm.Publish(commandTokens[1], strings.Join(commandTokens[2:], " "))
		manager.sendStringToClient("PUBACK", client)
	case "SUB":
		manager.qm.Subscribe(commandTokens[1], client)
	case "PINGREQ":
		manager.sendStringToClient("PINGRESP", client)
	default:
		manager.sendStringToClient(UNRECOGNISEDCOMMANDTEXT, client)
	}
}

func (manager *ConnectionManager) sendStringToClient(toSend string, client *Client) {
	fmt.Fprintln(client.Writer, toSend)
	client.Writer.Flush()
}
