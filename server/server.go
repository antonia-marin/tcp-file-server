package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	channels map[string]*channel
}

func (s *server) handleConnection(conn net.Conn) {
	for {
		b := make([]byte, 30112)
		_, err := conn.Read(b)
		if err != nil {
			log.Printf("Unable accept the request: %s", err.Error())
			return
		}

		cmd, chann, arg, cont := bytesParse(b)

		switch cmd {
		case "subscribe":
			s.subscribe(conn, chann, arg)
		case "channels":
			s.listChannels(conn)
		case "send":
			s.send(conn, chann, arg, cont)
		case "quit":
			s.quit(conn, arg)
		}
	}
}

func bytesParse(b []byte) (string, string, string, []byte) {
	cmdN := bytes.Index(b[:16], []byte{0})
	channN := bytes.Index(b[16:48], []byte{0})
	argsN := bytes.Index(b[48:112], []byte{0})

	cmd := string(b[:cmdN])
	channel := string(b[16:48][:channN])
	args := string(b[48:112][:argsN])
	content := b[112:]

	return cmd, channel, args, content
}

func (s *server) subscribe(c net.Conn, cName string, arg string) {
	chann, ok := s.channels[cName]
	if !ok {
		chann = &channel{
			name:    cName,
			members: map[string]net.Conn{c.RemoteAddr().String(): c},
		}
		s.channels[cName] = chann
	}

	chann.members[c.RemoteAddr().String()] = c

	if arg != "" {
		s.quitCurrentChannel(arg, c)
	}

	chann.messageBroadcast(c, fmt.Sprintf("New client has joined the channel"), s)
	s.msg(fmt.Sprintf("Welcome to %s", chann.name), c)
}

func (s *server) listChannels(c net.Conn) {
	var channels []string
	for name := range s.channels {
		channels = append(channels, name)
	}

	s.msg(fmt.Sprintf("Available channels are: %s", strings.Join(channels, ", ")), c)
}

func (s *server) send(c net.Conn, cName string, arg string, cont []byte) {
	channel, ok := s.channels[cName]
	if !ok {
		s.msg("You must join the room first", c)
		return
	}

	channel.fileBroadcast(c, arg, cont, s)
}

func (s *server) quit(c net.Conn, arg string) {
	log.Printf("Client has disconnected: %s", c.RemoteAddr().String())
	s.quitCurrentChannel(arg, c)

	go s.msg("Sad to see you go :(", c)
	c.Close()
}

func (s *server) msg(msg string, c net.Conn) {
	contenTypeBytes := make([]byte, 16)
	copy(contenTypeBytes, "message")

	extByte := make([]byte, 64)
	requestByte := append(contenTypeBytes, extByte...)

	msgBytes := make([]byte, 2048)
	copy(msgBytes, msg)
	requestByte = append(requestByte, msgBytes...)

	s.sendChannels(requestByte, c)
}

func (s *server) file(fileN string, cont []byte, c net.Conn) {
	contenTypeBytes := make([]byte, 16)
	copy(contenTypeBytes, "file")

	fileNameBytes := make([]byte, 64)
	copy(fileNameBytes, fileN)
	requestByte := append(contenTypeBytes, fileNameBytes...)
	requestByte = append(requestByte, cont...)

	s.sendChannels(requestByte, c)
}

func (s *server) sendChannels(responseByte []byte, c net.Conn) {
	_, err := c.Write(responseByte)
	if err != nil {
		fmt.Println("Server: failed to send content to client!")
	}
}

func (s *server) quitCurrentChannel(channelName string, c net.Conn) {
	channel := s.channels[channelName]
	if channel != nil {
		delete(channel.members, c.RemoteAddr().String())
		channel.messageBroadcast(c, fmt.Sprintf("A client has left the channel %s", channelName), s)
	}
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

	for {
		connection, err := ln.Accept()

		if err != nil {
			log.Printf("Unable accept connection: %s", err.Error())
			continue
		}

		log.Printf("New client has connected: %s", connection.RemoteAddr().String())
		s.msg(fmt.Sprintf("Welcome to the server client: %s", connection.RemoteAddr().String()), connection)
		go s.handleConnection(connection)
	}
}
