package meanstoend

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"runtime/debug"
	"strconv"
)

const REQUEST_LENGTH = 9

type store struct {
	rawData map[int32]int32
}

func (s *store) add(k, v int32) {
	s.rawData[k] = v
}
func (s *store) avg(k1, k2 int32) int32 {
	if k1 > k2 {
		return s.avg(k2, k1)
	}

	var sum int32 = 0
	var count int32 = 0

	for k := range s.rawData {
		if k >= k1 && k <= k2 {
			sum += s.rawData[k]
			count += 1
		}
	}

	if count == 0 {
		return 0
	}

	return sum / count
}

func assert(condition bool, msg string) {
	if !condition {
		debug.PrintStack()
		log.Fatalf("Assertion Error: %s", msg)
	}
}

type Request interface {
	opType() rune
	response(store) []byte
}

type request struct {
	raw []byte
}

func (r *request) opType() rune {
	assert(len(r.raw) == REQUEST_LENGTH, "Invalid request received")
	op := rune(r.raw[0])
	assert(op == 'I' || op == 'Q', "Invalid operation received")
	return op
}

func (r *request) response() []byte {
	panic("unimplemented")
}

type insertRequest struct {
	request
	timestamp int32
	price     int32
}

func (i *insertRequest) opType() rune {
	assert(i.request.opType() == 'I', "Invalid insert request opType")
	return 'I'
}

func (i *insertRequest) response(s store) []byte {
	s.add(i.timestamp, i.price)
	return []byte{'\r', '\n'}
}

type queryRequest struct {
	request
	timestampStart int32
	timestampEnd   int32
}

func (q *queryRequest) response(s store) []byte {
	avgPrice := s.avg(q.timestampStart, q.timestampEnd)
	var buffer []byte
	binary.BigEndian.PutUint32(buffer, uint32(avgPrice))
	buffer = append(buffer, '\r', '\n')
	return buffer
}

func (q *queryRequest) opType() rune {
	assert(q.request.opType() == 'Q', "Invalid query request opType")
	return 'Q'
}

func NewRequest(data []byte) Request {
	req := request{raw: data}
	op := req.opType()

	switch op {
	case 'I':
		return &insertRequest{
			request:   req,
			timestamp: int32(binary.BigEndian.Uint32(req.raw[1:4])),
			price:     int32(binary.BigEndian.Uint32(req.raw[4:8])),
		}
	case 'Q':
		return &queryRequest{
			request:        req,
			timestampStart: int32(binary.BigEndian.Uint32(req.raw[1:4])),
			timestampEnd:   int32(binary.BigEndian.Uint32(req.raw[4:8])),
		}

	default:
		log.Fatalf("Invalid request received")
	}

	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var s store

	for {
		buf := make([]byte, REQUEST_LENGTH)
		n, err := io.ReadFull(conn, buf)
		if err != nil {
			if err == io.EOF {
				log.Println("Connection closed by client")
				break
			}
			log.Printf("Error reading from connection: %s", err)
			break
		}
		assert(n == REQUEST_LENGTH, fmt.Sprintf("Did not read %d bytes", REQUEST_LENGTH))

		log.Printf("Recieved %v", strconv.Quote(string(buf)))

		request := NewRequest(buf)
		response := request.response(s)

		_, err = conn.Write(response)
		if err != nil {
			log.Printf("Something went wrong writing to connection: %s", err)
		}
	}
}

func MeanToEnd() {
	log.Println("Executing MeanToEnd")
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
