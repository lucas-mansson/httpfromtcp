package main

import (
	"fmt"
	"io"
	"net"
)

func getLinesChannel(in io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer in.Close()
		defer close(out)

		var currentLine []byte
		for {
			data := make([]byte, 8) // read chunks of up to 8 bytes
			nbrBytes, err := in.Read(data)

			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}

			currentSlice := data[:nbrBytes]

			for _, b := range currentSlice {
				if b == '\n' {
					out <- string(currentLine) // convert full UTF-8 line
					currentLine = nil
				} else {
					currentLine = append(currentLine, b)
				}
			}
		}

		// flush last line if not empty
		if len(currentLine) > 0 {
			out <- string(currentLine)
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

