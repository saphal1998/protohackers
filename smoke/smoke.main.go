package smoke

import (
	"io"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	data := make([]byte, 1024)
	output_data := make([]byte, 0)
	for {
		n, err := conn.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Something went wrong reading from connection: %s", err)
			break
		}
		output_data = append(output_data, data[:n]...)
	}

	_, err := conn.Write(output_data)
	if err != nil {
		log.Printf("Something went wrong writing to connection: %s", err)
	}
	return
}

func Smoke() {
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
