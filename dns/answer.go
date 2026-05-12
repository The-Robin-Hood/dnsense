package dns

import "encoding/binary"

// RR is a DNS Resource Record (an answer)
// RFC 1035 section 3.2.1
type RR struct {
	Name  string
	Type  uint16
	Class uint16
	TTL   uint32 // Time to live in seconds
	RData []byte // Actual record data, varies by type
}

// Pack serializes an RR to wire format
func (rr *RR) Pack() []byte {
	var buf []byte
	buf = append(buf, PackDomain(rr.Name)...)

	// Type, Class, TTL, RDLength, RData
	tmp := make([]byte, 8)
	binary.BigEndian.PutUint16(tmp[0:2], rr.Type)
	binary.BigEndian.PutUint16(tmp[2:4], rr.Class)
	binary.BigEndian.PutUint32(tmp[4:8], rr.TTL)
	buf = append(buf, tmp...)

	// RDLength + RData
	rdLen := make([]byte, 2)
	binary.BigEndian.PutUint16(rdLen, uint16(len(rr.RData)))
	buf = append(buf, rdLen...)
	buf = append(buf, rr.RData...)
	return buf
}
