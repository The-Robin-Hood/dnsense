package main

import (
	"dnsense/dns"
	"dnsense/handlers"
	"log"
	"os"
)

func main() {
	port := os.Getenv("DNS_PORT")
	if port == "" {
		port = "8053"
	}
	server := dns.NewServer("0.0.0.0:"+port, func(req *dns.Message) *dns.Message {
		return handlers.Route(req)
	})

	log.Fatal(server.Start())
}
