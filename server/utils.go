package main

import "bytes"

func BytesParse(b []byte) (string, string, string, []byte) {
	cmdN := bytes.Index(b[:16], []byte{0})
	channN := bytes.Index(b[16:48], []byte{0})
	argsN := bytes.Index(b[48:112], []byte{0})

	cmd := string(b[:cmdN])
	channel := string(b[16:48][:channN])
	args := string(b[48:112][:argsN])
	content := b[112:]

	return cmd, channel, args, content
}
