package main

import (
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	data := make([]byte, 1024)
	for {
		_, err := conn.Read(data)
		if err != nil {
			log.Printf("Something went wrong reading from connection: %s", err)
			break
		}
		_, err = conn.Write(data)
		if err != nil {
			log.Printf("Something went wrong writing to connection: %s", err)
			break
		}
	}

}

func main() {
	listener, err := net.Listen("tcp", ":5001")

	if err != nil {
		log.Fatalf("Something went wrong %s", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Could not accept connection %s", err)
		}
		go handleConnection(conn)
	}
}
