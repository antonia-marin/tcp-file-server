package main

import "net"

type channel struct {
	name    string
	members map[string]net.Conn
}

func (r *channel) messageBroadcast(conn net.Conn, msg string, s *server) {
	for addr, c := range r.members {
		if addr != conn.RemoteAddr().String() {
			s.msg(msg, c)
		}
	}
}

func (r *channel) fileBroadcast(conn net.Conn, fileN string, cont []byte, s *server) {
	for addr, c := range r.members {
		if addr != conn.RemoteAddr().String() {
			s.file(fileN, cont, c)
		}
	}
}
