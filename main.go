package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var inputFilePath = "messages.txt"

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("could not open file %s: %s\n", inputFilePath, err)
	}
	defer file.Close()

	buffer := make([]byte, 8)
	for {
		_, err := file.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Printf("Error while reading file: %v", err)
		}
		fmt.Printf("read: %s\n", string(buffer))
	}
}
