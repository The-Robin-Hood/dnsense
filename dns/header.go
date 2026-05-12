package dns

import (
	"encoding/binary"
	"fmt"
)

/**
Every DNS packet has this structure:

+---------------------+
|       Header        |  12 bytes, always
+---------------------+
|      Question       |  variable
+---------------------+
|      Answer         |  variable (in responses)
+---------------------+
|     Authority       |  variable (we'll get to this)
+---------------------+
|     Additional      |  variable
+---------------------+

1byte - 8bits
Here its 12 bytes - 96 bits

Which inturns split into 16bit per section so totally its 6 section as follows:

0  1  2  3  4  5  6  7  8  9  10 11 12 13 14 15
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                      ID                       |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    QDCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    ANCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    NSCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                    ARCOUNT                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

ID - 16 Random ID, clients sets and server echoes
QR - QUERY/RESPONSE 0/1
Opcode - most case its 0 means 'QUERY'
AA - AUTHORATIVE ANSWER - determine whether its from the owning zone source or from recursive resolver
TC - Truncated - response is way big so try with TCP (0/1)
RD - Recurion Desired - client asking server to recurse and give full resolution (0/1)
RA - Recursion Available - whether server can recurse or not (0/1)
RCODE - Response Code - 0 → OK, 3 → NXDOMAIN (domain doesn't exist), 2 → Server failure, 5 → Refused
QDCOUNT,ANCOUNT,NSCOUNT,ARCOUNT - Number of Questions,Answers,Authority Records and Additional Records

**/

type Header struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

const (
	FlagQR     = 1 << 15
	FlagAA     = 1 << 10
	FlagTC     = 1 << 9
	FlagRD     = 1 << 8
	FlagRA     = 1 << 7
	RCodeMask  = 0x000F // Bottom 4 bits
	OpcodeMask = 0x7800 // Bits 11-14
)

func (h *Header) IsItResponse() bool {
	return h.Flags&FlagQR != 0
}

func (h *Header) SetResponse() {
	h.Flags |= FlagQR
}

func (h *Header) SetAA() {
	h.Flags |= FlagAA
}

func (h *Header) SetRD() {
	h.Flags |= FlagRD
}

func (h *Header) SetRA() {
	h.Flags |= FlagRA
}

func (h *Header) SetRCode(code uint8) {
	h.Flags = (h.Flags &^ RCodeMask) | uint16(code)
}

func (h *Header) GetRCode() uint8 {
	return uint8(h.Flags & RCodeMask)
}

// Pack serializes the header to 12 bytes (big-endian per RFC)
func (h *Header) Pack() []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint16(buf[0:2], h.ID)
	binary.BigEndian.PutUint16(buf[2:4], h.Flags)
	binary.BigEndian.PutUint16(buf[4:6], h.QDCount)
	binary.BigEndian.PutUint16(buf[6:8], h.ANCount)
	binary.BigEndian.PutUint16(buf[8:10], h.NSCount)
	binary.BigEndian.PutUint16(buf[10:12], h.ARCount)
	return buf
}

// UnpackHeader parses a 12-byte DNS header
func UnpackHeader(buf []byte) (*Header, error) {
	if len(buf) < 12 {
		return nil, fmt.Errorf("buffer too short for header: %d bytes", len(buf))
	}
	return &Header{
		ID:      binary.BigEndian.Uint16(buf[0:2]),
		Flags:   binary.BigEndian.Uint16(buf[2:4]),
		QDCount: binary.BigEndian.Uint16(buf[4:6]),
		ANCount: binary.BigEndian.Uint16(buf[6:8]),
		NSCount: binary.BigEndian.Uint16(buf[8:10]),
		ARCount: binary.BigEndian.Uint16(buf[10:12]),
	}, nil
}
