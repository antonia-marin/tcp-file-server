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

func (c *client) readInput() {
	for {
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Fatalf("Error: %s", err.Error())
			return
		}

		msg := strings.Trim(line, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/subscribe":
			c.subscribe(args)
		case "/channels":
			c.channels()
		case "/send":
			c.send(args)
		case "/quit":
			c.quit()
		default:
			fmt.Println("Error Command not found: ", cmd)
		}
	}
}

func (c *client) subscribe(arguments []string) {}

func (c *client) channels() {}

func (c *client) send(arguments []string) {}

func (c *client) quit() {}

func main() {
	conn, err := net.Dial("tcp", ":9999")

	if err != nil {
		log.Fatalf("Unable to connect with the server: %s", err.Error())
		return
	}

	c := &client{
		conn: conn,
	}

	c.readInput()
}
