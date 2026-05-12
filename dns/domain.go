package dns

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// PackDomain encodes a domain name into DNS wire format
// "ansari.wtf" → [6]ansari[3]wtf[0]
func PackDomain(domain string) []byte {
	// Remove trailing dot if present
	domain = strings.TrimSuffix(domain, ".")

	if domain == "" {
		return []byte{0} // root
	}

	var buf []byte
	labels := strings.SplitSeq(domain, ".")
	for label := range labels {
		if len(label) > 63 {
			panic("label too long") // RFC limit: 63 chars per label
		}
		buf = append(buf, byte(len(label)))
		buf = append(buf, []byte(label)...)
	}
	buf = append(buf, 0) // null terminator
	return buf
}

// UnpackDomain decodes a DNS name from wire format
// Returns the name and how many bytes were consumed
// Also handles compression pointers (0xC0 prefix)
func UnpackDomain(buf []byte, offset int) (string, int, error) {
	var labels []string
	visited := make(map[int]bool) // detect pointer loops
	origOffset := offset
	jumped := false
	jumpOffset := 0

	for {
		if offset >= len(buf) {
			return "", 0, fmt.Errorf("name parse out of bounds at offset %d", offset)
		}

		length := int(buf[offset])

		// Check for compression pointer: top 2 bits are 11 (0xC0)
		if length&0xC0 == 0xC0 {
			if offset+1 >= len(buf) {
				return "", 0, fmt.Errorf("compression pointer out of bounds")
			}
			// Pointer is 14-bit offset into the message
			ptr := int(binary.BigEndian.Uint16(buf[offset:offset+2]) &^ 0xC000)

			if visited[ptr] {
				return "", 0, fmt.Errorf("compression pointer loop detected")
			}
			visited[ptr] = true

			if !jumped {
				jumpOffset = offset + 2 // after the pointer, this is where we resume
			}
			jumped = true
			offset = ptr
			continue
		}

		// Normal label
		if length == 0 {
			// End of name
			offset++ // consume the null byte
			break
		}

		offset++ // move past the length byte
		if offset+length > len(buf) {
			return "", 0, fmt.Errorf("label out of bounds")
		}
		labels = append(labels, string(buf[offset:offset+length]))
		offset += length
	}

	name := strings.Join(labels, ".") + "."

	// If we jumped, return the offset after the pointer (2 bytes)
	// If we didn't jump, return where we ended up
	consumed := offset - origOffset
	if jumped {
		consumed = jumpOffset - origOffset
	}

	return name, consumed, nil
}

