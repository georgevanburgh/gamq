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
	UNRECOGNISEDCOMMANDTEXT = "Unrecognised command"
	HELPSTRING              = `Available commands:
	HELP: Prints this text
	PUB <queue> <message>: Publish <message> to <queue>
	SUB <queue>: Subscribe to messages on <queue>
	PINGREQ: Requests a PINGRESP from the server`
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

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", Configuration.Port))

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
		log.Debug("A new connection was opened.")
		manager.wg.Add(1)
		go manager.handleConnection(&conn)
	}
}

func (manager *ConnectionManager) handleConnection(conn *net.Conn) {
	defer manager.wg.Done()
	connReader := bufio.NewReader(*conn)
	connWriter := bufio.NewWriter(*conn)
	closedChannel := make(chan bool)
	client := Client{Name: strconv.Itoa(manager.rand.Int()),
		Writer: connWriter,
		Reader: connReader,
		Closed: &closedChannel}

	for {
		// Read until newline
		line, err := client.Reader.ReadString('\n')

		if err != nil {
			// Connection has been closed
			log.Debugf("%s closed connection", client.Name)
			*client.Closed <- true
			break
		}

		// Parse command and (optionally) return response (if any)
		manager.parseClientCommand(line, &client)
	}

	log.Info("A connection was closed")
}

func (manager *ConnectionManager) parseClientCommand(commandLine string, client *Client) {

	commandTokens := strings.Fields(string(commandLine[:len(commandLine)]))

	if len(commandTokens) == 0 {
		return
	}

	commandTokens[0] = strings.ToUpper(commandTokens[0])

	switch commandTokens[0] {
	case "HELP":
		manager.sendStringToClient(HELPSTRING, client)
	case "PUB":
		message := strings.Join(commandTokens[2:], " ")
		manager.qm.Publish(commandTokens[1], &message)
		// manager.sendStringToClient("PUBACK", client)
	case "SUB":
		manager.qm.Subscribe(commandTokens[1], client)
	case "PINGREQ":
		manager.sendStringToClient("PINGRESP", client)
	case "CLOSE":
		manager.qm.CloseQueue(commandTokens[1])
	default:
		manager.sendStringToClient(UNRECOGNISEDCOMMANDTEXT, client)
	}
}

func (manager *ConnectionManager) sendStringToClient(toSend string, client *Client) {
	fmt.Fprintln(client.Writer, toSend)
	client.Writer.Flush()
}
