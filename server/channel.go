package main

import "net"

type channel struct {
	name    string
	members map[string]net.Conn
}

func (chann *channel) messageBroadcast(conn net.Conn, msg string, s *server) {
	for addr, c := range chann.members {
		if addr != conn.RemoteAddr().String() {
			s.msg(msg, c)
		}
	}
}

func (chann *channel) fileBroadcast(conn net.Conn, fileN string, cont []byte, s *server) {
	for addr, c := range chann.members {
		if addr != conn.RemoteAddr().String() {
			s.file(fileN, cont, c)
		}
	}
}

func (chann *channel) channMembers() []string {
	var channelMembers []string
	for addr, _ := range chann.members {
		channelMembers = append(channelMembers, addr)
	}

	return channelMembers
}
