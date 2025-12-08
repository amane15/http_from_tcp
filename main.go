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

	linesChan := getLinesChannel(file)
	for line := range linesChan {
		fmt.Printf("read: %s\n", line)
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
