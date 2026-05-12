package dns

import (
	"fmt"
	"log"
	"net"
)

type Handler func(req *Message) *Message

type Server struct {
	addr    string
	conn    *net.UDPConn
	handler Handler
}

func NewServer(addr string, handler Handler) *Server {
	return &Server{
		addr:    addr,
		handler: handler,
	}
}

func (s *Server) Start() error {
	udpAddr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return fmt.Errorf("resolve addr: %w", err)
	}

	// Bind the UDP socket
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("listen udp: %w", err)
	}
	s.conn = conn

	log.Printf("DNS server listening on %s", s.addr)

	// Read loop — runs forever
	buf := make([]byte, 512) // RFC 1035: max UDP DNS message is 512 bytes
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("read error: %v", err)
			continue
		}

		// Handle each request in a goroutine
		// so slow requests don't block other clients
		go s.handlePacket(buf[:n], clientAddr)
	}
}

func (s *Server) handlePacket(buf []byte, clientAddr *net.UDPAddr) {
	// Parse the incoming message
	req, err := UnpackMessage(buf)
	if err != nil {
		log.Printf("failed to parse DNS message from %s: %v", clientAddr, err)
		return
	}

	log.Printf("query from %s: %s type=%d", clientAddr, req.Questions[0].Name, req.Questions[0].Type)

	// Call our handler to build a response
	resp := s.handler(req)
	if resp == nil {
		return
	}

	// Serialize response to wire format
	respBytes, err := resp.Pack()
	if err != nil {
		log.Printf("failed to pack response: %v", err)
		return
	}

	// Send it back to the client
	_, err = s.conn.WriteToUDP(respBytes, clientAddr)
	if err != nil {
		log.Printf("failed to send response to %s: %v", clientAddr, err)
	}
}
