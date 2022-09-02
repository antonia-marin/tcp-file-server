package main

import "net"

type channel struct {
	name    string
	members map[string]net.Conn
}
