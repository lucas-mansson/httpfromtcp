package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	
	file, err := os.Open("messages.txt")

	if err != nil {
		panic(err)
	}

	for {
		data := make([]byte, 8) // Allocate and initialize empty array of bytes
		nbrBytes, err := file.Read(data) // reads 8 bytes from file and stores in data
		
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		fmt.Printf("read: %s\n", string(data[:nbrBytes])) // returns and removes the first n bytes from the array
	}

}
