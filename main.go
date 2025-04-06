package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatalf("error while opening file: %v", err)
	}

	buffer := make([]byte, 8, 8)
	currentLine := ""

	for {
		n, err := file.Read(buffer)
		fmt.Println("Current Line ->", currentLine)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("error while reading data: %v", err)
			break
		}

		str := string(buffer[:n])
		parts := strings.Split(str, "\n")

		for i := 0; i < len(parts)-1; i++ {
			fmt.Println("read:", currentLine+parts[i])
			currentLine = ""
		}
		currentLine += parts[len(parts)-1]
	}
}
