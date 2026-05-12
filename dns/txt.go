package dns

import "fmt"

const (
	maxTXTString  = 255 // RFC 1035: each string max 255 bytes
	maxUDPPayload = 512 // RFC 1035: max UDP DNS message
)

// PackTXT encodes a slice of strings into TXT RData wire format
// Each string becomes a length-prefixed chunk
// "hello" "world" → 0x05 h e l l o 0x05 w o r l d
func PackTXT(strings []string) ([]byte, error) {
	var rdata []byte

	for _, s := range strings {
		if len(s) > maxTXTString {
			return nil, fmt.Errorf("TXT string exceeds 255 bytes: %d", len(s))
		}
		rdata = append(rdata, byte(len(s)))
		rdata = append(rdata, []byte(s)...) 
	}

	return rdata, nil
}

// ChunkText splits a long string into 255-byte chunks
// ready to be passed into PackTXT
func ChunkText(text string) []string {
	var chunks []string

	for len(text) > 0 {
		if len(text) <= maxTXTString {
			chunks = append(chunks, text)
			break
		}
		// Try to split on a word boundary near the 255 limit
		cutoff := maxTXTString
		for cutoff > 200 && text[cutoff] != ' ' {
			cutoff--
		}
		chunks = append(chunks, text[:cutoff])
		text = text[cutoff+1:] 
	}

	return chunks
}

// NewTXTRecord builds a complete RR with TXT type from a long string
// It handles chunking automatically
func NewTXTRecord(name string, ttl uint32, text string) (*RR, error) {
	chunks := ChunkText(text)

	rdata, err := PackTXT(chunks)
	if err != nil {
		return nil, fmt.Errorf("pack txt: %w", err)
	}

	return &RR{
		Name:  name,
		Type:  TypeTXT,
		Class: ClassIN,
		TTL:   ttl,
		RData: rdata,
	}, nil
}
