package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/websocket"
)

type server struct {
	channels map[string]*channel
	clients  map[string]net.Conn
	files    map[string][]string
}

type Statistics struct {
	Clients           []string            `json:"clients"`
	Channels          []string            `json:"channels"`
	ClientsOnChannels map[string][]string `json:"clientsOnChannels"`
	Files             map[string][]string `json:"files"`
}

func (s *server) handleConnection(conn net.Conn) {
	for {
		b := make([]byte, 60112)
		_, err := conn.Read(b)
		if err != nil {
			log.Printf("Unable accept the request: %s", err.Error())
			return
		}

		cmd, chann, arg, cont := BytesParse(b)

		switch cmd {
		case "subscribe":
			s.subscribe(conn, chann, arg)
		case "channels":
			s.listChannels(conn)
		case "send":
			s.send(conn, chann, arg, cont)
		case "quit":
			s.quit(conn, chann)
		}
	}
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
	s.files[c.RemoteAddr().String()] = append(s.files[c.RemoteAddr().String()], arg)
}

func (s *server) quit(c net.Conn, cName string) {
	s.quitCurrentChannel(cName, c)
	s.msg("Sad to see you go :(", c)

	log.Printf("Client has disconnected: %s", c.RemoteAddr().String())
	delete(s.clients, c.RemoteAddr().String())
	c.Close()
}

func (s *server) msg(msgContent string, c net.Conn) {
	contenTypeBytes := make([]byte, 16)
	copy(contenTypeBytes, "message")

	extByte := make([]byte, 64)
	requestByte := append(contenTypeBytes, extByte...)

	msgBytes := make([]byte, 2048)
	copy(msgBytes, msgContent)
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
		log.Printf("Error server: failed to send content to client!")
	}
}

func (s *server) quitCurrentChannel(channelName string, c net.Conn) {
	channel := s.channels[channelName]
	if channel != nil {
		delete(channel.members, c.RemoteAddr().String())
		channel.messageBroadcast(c, fmt.Sprintf("A client has left the channel %s", channelName), s)
	}
}

func (s *server) handleUIConnection(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error WS Upgrader: %v", err)
	}

	defer ws.Close()

	for {
		var msg string
		errRead := ws.ReadJSON(&msg)
		if errRead != nil {
			log.Printf("Error WS on ReadJSON: %v", errRead)
			break
		}

		stat := &Statistics{
			Clients:           s.serverClients(),
			Channels:          s.serverChannels(),
			ClientsOnChannels: buildMapClientsOnChannels(s.channels),
			Files:             s.files,
		}

		errWrite := ws.WriteJSON(stat)
		if errWrite != nil {
			log.Printf("Error WS on WriteJSON: %v", errWrite)
			ws.Close()
		}
	}
}

func (s *server) serverUI() {
	http.HandleFunc("/ws", s.handleUIConnection)
	err2 := http.ListenAndServe("localhost:8080", nil)
	if err2 != nil {
		log.Fatal("ListenAndServe: ", err2)
	}
}

func (s *server) serverClients() []string {
	clients := s.clients
	var keys []string
	for k := range clients {
		keys = append(keys, k)
	}

	return keys
}

func (s *server) serverChannels() []string {
	channels := s.channels
	var keys []string
	for k := range channels {
		keys = append(keys, k)
	}

	return keys
}

func buildMapClientsOnChannels(serverChannels map[string]*channel) map[string][]string {
	clientsOnChannels := make(map[string][]string)
	keysChannels := reflect.ValueOf(serverChannels).MapKeys()
	for _, ss := range keysChannels {
		channel, _ := serverChannels[ss.String()]
		chanMem := channel.channMembers()
		clientsOnChannels[ss.String()] = chanMem
	}

	return clientsOnChannels
}

func main() {
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatalf("Error: Unable to start server %s", err.Error())
	}

	defer ln.Close()

	s := &server{
		channels: make(map[string]*channel),
		clients:  make(map[string]net.Conn),
		files:    make(map[string][]string),
	}

	go s.serverUI()

	for {
		connection, err := ln.Accept()
		if err != nil {
			log.Printf("Error: Unable accept connection %s", err.Error())
			continue
		}

		log.Printf("New client has connected: %s", connection.RemoteAddr().String())
		s.msg(fmt.Sprintf("Welcome to the server client: %s", connection.RemoteAddr().String()), connection)
		s.clients[connection.RemoteAddr().String()] = connection
		go s.handleConnection(connection)
	}
}
