package main

import (
	"bufio"
	"fmt"
	"io"
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

func (c *client) buildSendPayload(arguments []string) []byte {
	byteFile := readFile(arguments[1])
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

func readFile(filePath string) []byte {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("OS. Stat() function execution error, error is:% v \n", err)
		return nil
	}

	fileNameBytes := make([]byte, 64)
	copy(fileNameBytes, fileInfo.Name())

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("OS. Open() function execution error, error is:% v \n", err)
		return nil
	}

	defer file.Close()

	fileBytes := make([]byte, 30112)
	n, err := file.Read(fileBytes)
	if err != nil {
		if err != io.EOF {
			fmt.Printf("file. Read() method execution error, error is:% v \n", err)
		}
		return nil
	}

	requestBytes := append(fileNameBytes, fileBytes[:n]...)
	return requestBytes
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
