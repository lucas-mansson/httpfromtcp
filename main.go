package main

import (
	"fmt"
	"io"
	"os"
)

func getLinesChannel(file io.ReadCloser) <- chan string {
	out := make(chan string, 1)

	go func() {
		defer file.Close()
		defer close(out)

		var currentLine string
		for {
			data := make([]byte, 8) // Allocate and initialize empty array of 8 bytes
			nbrBytes, err := file.Read(data) // reads 8 bytes from file and stores in data

			if err == io.EOF { break }
			if err != nil { panic(err) }

			currentSlice := data[:nbrBytes]

			for _, byte := range currentSlice {
				if byte == '\n' {
					out <- string(currentLine)
					currentLine = ""
				} else {
					currentLine += string(byte)
				}
			}
		}
	}()
	return out
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}
	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read %s\n", line)
	}
}

