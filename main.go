package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatalf("error while opening file: %v", err)
	}

	data := make([]byte, 8)

	for {
		n, err := file.Read(data)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("error while reading data: %v", err)
			break
		}

		str := string(data[:n])
		fmt.Printf("read: %s\n", str)

	}
}
