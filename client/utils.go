package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func ReadInput() ([]string, *string) {
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return nil, nil
	}

	msg := strings.Trim(line, "\r\n")
	args := strings.Split(msg, " ")
	cmd := strings.TrimSpace(args[0])

	return args, &cmd
}

func BytesParse(b []byte) (string, string, []byte) {
	typeN := bytes.Index(b[:16], []byte{0})
	fileN := bytes.Index(b[16:80], []byte{0})

	typeCont := string(b[:typeN])
	fileName := string(b[16:80][:fileN])
	content := b[80:]

	return typeCont, fileName, content
}

func ReadFile(filePath string) []byte {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("Error: OS. Stat() function execution error %s \n", err.Error())
		return nil
	}

	fileNameBytes := make([]byte, 64)
	copy(fileNameBytes, fileInfo.Name())

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error: OS. Open() function execution error %s \n", err.Error())
		return nil
	}

	defer file.Close()

	fileBytes := make([]byte, 60000)
	n, err := file.Read(fileBytes)
	if err != nil {
		if err != io.EOF {
			fmt.Printf("Error: file. Read() method execution error %s \n", err.Error())
		}
		return nil
	}

	requestBytes := append(fileNameBytes, fileBytes[:n]...)
	return requestBytes
}
