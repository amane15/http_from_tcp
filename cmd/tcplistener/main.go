package main

import (
	"fmt"
	"log"
	"net"

	"github.com/amane15/http_from_tcp/internal/request"
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

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error while reading request: %s\n", err.Error())
			return
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		conn.Close()
		fmt.Println("channel is closed, closing connection")
	}
}
