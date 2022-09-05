package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type client struct {
	conn    net.Conn
	channel string
}

func readInput() ([]string, *string) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')

	if err != nil {
		log.Fatalf("Error: %s", err.Error())
		return nil, nil
	}

	msg := strings.Trim(line, "\r\n")
	args := strings.Split(msg, " ")
	cmd := strings.TrimSpace(args[0])

	return args, &cmd
}

func (c *client) handleInput() {
	for {
		args, cmd := readInput()
		var bytesRequest []byte

		switch *cmd {
		case "/subscribe":
			bytesRequest = c.buildSubscribePayload(args)
		case "/channels":
			bytesRequest = buildChannelsPayload()
		case "/send":
			bytesRequest = c.buildSendPayload(args)
		case "/quit":
			bytesRequest = buildQuitPayload()
		default:
			fmt.Println("Error Command not found: ", *cmd)
		}

		c.request(bytesRequest)
	}
}

func (c *client) buildSubscribePayload(arguments []string) []byte {
	commBytes := make([]byte, 16)
	copy(commBytes, "subscribe")

	clientChannel := []string{strings.Join(arguments[1:], "-")}
	channelBytes := make([]byte, 32)
	copy(channelBytes, clientChannel[0])

	requestByte := append(commBytes, channelBytes...)

	if c.channel != "" {
		oldChannelNameBytes := make([]byte, 64)
		copy(oldChannelNameBytes, c.channel)
		requestByte = append(requestByte, oldChannelNameBytes...)
	}
	c.channel = clientChannel[0]

	return requestByte
}

func buildChannelsPayload() []byte {
	commBytes := make([]byte, 16)
	copy(commBytes, "channels")

	return commBytes
}

func (c *client) buildSendPayload(arguments []string) []byte { return nil }

func buildQuitPayload() []byte {
	commBytes := make([]byte, 16)
	copy(commBytes, "quit")

	return commBytes
}

func (c *client) request(request []byte) {
	if request != nil && len(request) > 0 {
		c.conn.Write(request)
	}
}

func main() {
	conn, err := net.Dial("tcp", ":9999")

	if err != nil {
		log.Fatalf("Unable to connect with the server: %s", err.Error())
		return
	}

	c := &client{
		conn: conn,
	}

	c.handleInput()
}
