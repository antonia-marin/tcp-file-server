package main

import (
	"log"
	"net"
)

type server struct {
	channels map[string]*channel
}

func main() {
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatalf("Unable to start server: %s", err.Error())
	}

	defer ln.Close()

	s := &server{
		channels: make(map[string]*channel),
	}
}
