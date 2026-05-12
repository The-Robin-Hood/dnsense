package dns

import "fmt"

type Message struct {
	Header    Header
	Questions []Question
	Answers   []RR
}

// Pack serializes the entire DNS message to wire format
func (m *Message) Pack() ([]byte, error) {
	var buf []byte

	m.Header.QDCount = uint16(len(m.Questions))
	m.Header.ANCount = uint16(len(m.Answers))
	m.Header.NSCount = 0
	m.Header.ARCount = 0

	buf = append(buf, m.Header.Pack()...)

	for _, q := range m.Questions {
		buf = append(buf, q.Pack()...)
	}

	for _, rr := range m.Answers {
		buf = append(buf, rr.Pack()...)
	}

	return buf, nil
}

// Unpack parses a raw DNS message from wire format
func UnpackMessage(buf []byte) (*Message, error) {
	if len(buf) < 12 {
		return nil, fmt.Errorf("message too short: %d bytes", len(buf))
	}

	header, err := UnpackHeader(buf)
	if err != nil {
		return nil, err
	}

	msg := &Message{Header: *header}
	offset := 12 // start after header

	// Parse questions
	for i := 0; i < int(header.QDCount); i++ {
		q, consumed, err := UnpackQuestion(buf, offset)
		if err != nil {
			return nil, fmt.Errorf("question %d: %w", i, err)
		}
		msg.Questions = append(msg.Questions, *q)
		offset += consumed
	}

	return msg, nil
}
