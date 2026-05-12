package dns

import (
	"encoding/binary"
	"fmt"
)

// Record types per RFC 1035 + RFC 3596
const (
	TypeA     = 1
	TypeNS    = 2
	TypeCNAME = 5
	TypeMX    = 15
	TypeTXT   = 16
	TypeAAAA  = 28 // RFC 3596

	ClassIN = 1 // Internet class
)

type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

// Pack serializes a question to wire format
func (q *Question) Pack() []byte {
	var buf []byte
	buf = append(buf, PackDomain(q.Name)...)
	buf = append(buf, 0, 0) // Type
	buf = append(buf, 0, 0) // Class
	binary.BigEndian.PutUint16(buf[len(buf)-4:], q.Type)
	binary.BigEndian.PutUint16(buf[len(buf)-2:], q.Class)
	return buf
}

// UnpackQuestion parses a question from wire format at offset
func UnpackQuestion(buf []byte, offset int) (*Question, int, error) {
	name, consumed, err := UnpackDomain(buf, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("question name: %w", err)
	}
	offset += consumed

	if offset+4 > len(buf) {
		return nil, 0, fmt.Errorf("question too short for type/class")
	}

	qtype := binary.BigEndian.Uint16(buf[offset : offset+2])
	qclass := binary.BigEndian.Uint16(buf[offset+2 : offset+4])

	return &Question{
		Name:  name,
		Type:  qtype,
		Class: qclass,
	}, consumed + 4, nil
}
