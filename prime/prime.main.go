package prime

import (
	"bufio"
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
	log.Printf("Got connection: %v", conn.RemoteAddr())
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error reading from connection: %v", err)
			return
		}

		var req request
		err = json.Unmarshal([]byte(message), &req)
		if err != nil || req.Method != "isPrime" {
			malformedResponse := []byte(`{"error": "malformed request"}\n`)
			conn.Write(malformedResponse)
			return
		}

		number_is_prime := checkPrime(req.Number)
		correct_response := response{Method: req.Method,
			Prime: number_is_prime,
		}

		responseBytes, err := json.Marshal(correct_response)
		if err != nil {
			log.Printf("Error marshalling response: %v", err)
			return
		}
		log.Printf("Sending back: %v for %v", correct_response, req)

		conn.Write(append(responseBytes, '\n'))
	}
}

func checkPrime(f float64) bool {
	if f != math.Floor(f) {
		return false
	}
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
	for i := 3; i <= limit; i += 2 {
		if number%i == 0 {
			return false
		}
	}

	return true
}
