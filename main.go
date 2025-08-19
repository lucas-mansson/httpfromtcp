package main

import (
	"fmt"
	"io"
	"net"
)

func getLinesChannel(in io.ReadCloser) <- chan string {
	out := make(chan string, 1)

	go func() {
		defer in.Close()
		defer close(out)

		var currentLine string
		for {
			data := make([]byte, 8) // Allocate and initialize empty array of 8 bytes
			nbrBytes, err := in.Read(data) // reads 8 bytes from file and stores in data

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
	port := ":42069"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	defer fmt.Print("Connection closed.")
	
	fmt.Printf("Listening on port %s\n", port)
	for {
		conn, err := listener.Accept()
		fmt.Print("Connection accepted!")
		if err != nil {
			panic(err)
		}
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("read %s\n", line)
		}
	}

}

