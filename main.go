package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	tcpListener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("error creating tcp listener: %s\n", err.Error())
	}
	defer tcpListener.Close()

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Fatalf("error accepting connection: %s\n", err.Error())
		}
		fmt.Printf("connection at address %s has been established\n", conn.LocalAddr().String())
		linesChan := getLinesChannel(conn)
		for line := range linesChan {
			fmt.Println(line)
		}
		conn.Close()
		fmt.Println("channel is closed, closing connection")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer f.Close()
		defer close(lines)
		buffer := make([]byte, 8)
		var currentLine []byte

		for {
			n, err := f.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatalf("error while reading file: %s", err)
				break
			}
			if n > 0 {
				chunk := buffer[:n]
				for _, b := range chunk {
					if b == '\n' {
						// fmt.Printf("read: %s\n", string(currentLine))
						lines <- string(currentLine)
						currentLine = currentLine[:0]
					} else {
						currentLine = append(currentLine, b)
					}
				}
			}
		}
		if len(currentLine) > 0 {
			// fmt.Printf("read: %s\n", string(currentLine))
			lines <- string(currentLine)
		}
	}()

	return lines
}
