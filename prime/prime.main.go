package prime

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"net"
)

type request struct {
	Method string  `json:"method"`
	Number float64 `json:"number"`
}

type response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func Prime() {
	log.Println("Executing PRIME")
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

	log.Printf("Received %v", string(output_data))
	// We have the output data
	var req request
	err := json.Unmarshal(output_data, &req)
	if err != nil || req.Method != "isPrime" {
		if err != nil {
			log.Printf("Something went wrong when reading the payload: %s", err)
		}
		_, err = conn.Write([]byte{'\r', '\n'})
		if err != nil {
			log.Printf("Something went wrong writing to connection: %s", err)
			return
		}
	}

	log.Printf("Received as object: %v", req)
	number_is_prime := checkPrime(req.Number)
	correct_response := response{
		Method: "isPrime",
		Prime:  number_is_prime,
	}

	log.Printf("Sending object: %v", correct_response)
	response, err := json.Marshal(correct_response)
	if err != nil {
		log.Printf("Something went wrong when marshalling response: %s", err)
		return
	}
	_, err = conn.Write(response)
	if err != nil {
		log.Printf("Something went wrong writing correct response to connection: %s", err)
		return
	}
	return
}

func checkPrime(f float64) bool {
	number := int(f)

	if number <= 1 {
		return false
	}

	if number == 2 {
		return true
	}

	if number%2 == 0 {
		return false
	}

	limit := int(math.Floor(math.Sqrt(float64(number))))
	for i := 3; i <= limit; i += 1 {
		if number%i == 0 {
			return false
		}
	}

	return true
}
