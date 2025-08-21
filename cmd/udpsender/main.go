package main

import (
	"os"
	"fmt"
	"net"
	"bufio"
)

func main() {
	recieverAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, recieverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		in, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print(err)
		}
		conn.Write([]byte(in))
	}
}
