package main

import (
	"dnsense/dns"
	"dnsense/handlers"
	"log"
)

func main() {
	server := dns.NewServer("0.0.0.0:8053", func(req *dns.Message) *dns.Message {
		return handlers.Route(req)
	})

	log.Println("dns running on :8053")
	log.Fatal(server.Start())
}
