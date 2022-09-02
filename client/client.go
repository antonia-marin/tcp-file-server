package main

import (
	"log"
	"net"
)

type client struct {
	conn    net.Conn
	channel string
}

func newClient(conn net.Conn) {
	c := &client{
		conn: conn,
	}
}

func main() {
	conn, err := net.Dial("tcp", ":9999")

	if err != nil {
		log.Fatalf("Unable to connect with the server: %s", err.Error())
		return
	}

	newClient(conn)
}
