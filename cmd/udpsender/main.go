package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		log.Fatalf("error resolving udp address: %s\n", err.Error())
	}

	conn, err := net.DialUDP("udp", udpAddr, udpAddr)
	if err != nil {
		log.Fatalf("error creating udp connection: %s\n", err.Error())
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">")
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("error reading input from stdin: %s\n", err.Error())
		}
		_, err = conn.Write([]byte(input))
		if err != nil {
			log.Fatalf("error writing to a udp connection: %s\n", err.Error())
		}
	}
}
