package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func getLinesChannel(in io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer in.Close()
		defer close(out)

var currentLine string
	for {
		data := make([]byte, 8) // Allocate and initialize empty array of bytes
		nbrBytes, err := in.Read(data) // reads 8 bytes from file and stores in data
		if err == io.EOF { break }
		if err != nil { panic(err) }
		
		data = data[:nbrBytes]
		if i := bytes.IndexByte(data, '\n'); i != -1 {
			currentLine += string(data[:i])
			data = data[i+1:]
			fmt.Printf("read: %s\n", currentLine) 
			currentLine = ""
		}
		currentLine += string(data)
	}

	if len(currentLine) != 0 {
		fmt.Printf("read: %s\n", currentLine) 
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
		fmt.Println("Connection accepted")
		if err != nil {
			panic(err)
		}
		for line := range getLinesChannel(conn) {
			fmt.Printf(line)
		}
	}

}

