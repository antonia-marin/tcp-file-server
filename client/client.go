package main

import (
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

func (c *client) handleInput() {
	for {
		args, cmd := ReadInput()
		var bytesRequest []byte

		switch *cmd {
		case "/subscribe":
			bytesRequest = c.buildSubscribePayload(args)
		case "/channels":
			bytesRequest = buildChannelsPayload()
		case "/send":
			bytesRequest = c.buildSendPayload(args)
		case "/quit":
			bytesRequest = c.buildQuitPayload()
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

func (c *client) buildSendPayload(arguments []string) []byte {
	byteFile := ReadFile(arguments[1])
	if byteFile == nil {
		return nil
	}

	commBytes := make([]byte, 16)
	copy(commBytes, "send")

	channelBytes := make([]byte, 32)
	copy(channelBytes, c.channel)
	requestBytes := append(commBytes, channelBytes...)
	requestBytes = append(requestBytes, byteFile...)

	return requestBytes
}

func (c *client) buildQuitPayload() []byte {
	commBytes := make([]byte, 16)
	copy(commBytes, "quit")

	channelBytes := make([]byte, 32)
	copy(channelBytes, c.channel)

	requestBytes := append(commBytes, channelBytes...)

	return requestBytes
}

func (c *client) request(req []byte) {
	if req != nil && len(req) > 0 {
		c.conn.Write(req)
	}
}

func handleResponse(conn net.Conn) {
	for {
		b := make([]byte, 60112)
		_, err := conn.Read(b)
		if err != nil {
			log.Printf("Error: Unable accept the request %s", err.Error())
			return
		}

		typeCont, fileName, content := BytesParse(b)
		switch typeCont {
		case "message":
			msg(content)
		case "file":
			file(fileName, content)
		}
	}
}

func msg(msgContent []byte) {
	fmt.Printf("> %s \n", string(msgContent))
}

func file(fileName string, cont []byte) {
	newFile, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error: OS. Create() function execution error %s \n", err.Error())
		return
	}

	newFile.Write(cont)
}

func main() {
	conn, err := net.Dial("tcp", ":9999")

	if err != nil {
		log.Fatalf("Error: Unable to connect with the server %s", err.Error())
		return
	}

	c := &client{
		conn: conn,
	}

	go handleResponse(conn)
	c.handleInput()
}
