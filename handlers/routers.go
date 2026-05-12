package handlers

import (
	"dnsense/dns"
	"strings"
)

func baseResponse(req *dns.Message) *dns.Message {
	resp := &dns.Message{}
	resp.Header = req.Header
	resp.Header.SetResponse()
	resp.Header.SetAA()
	resp.Header.Flags &^= (1 << 5) // clear AD bit
	resp.Header.ARCount = 0
	resp.Questions = req.Questions
	return resp
}

func errorResponse(req *dns.Message, msg string) *dns.Message {
	resp := baseResponse(req)

	rr, err := dns.NewTXTRecord(req.Questions[0].Name, 30, "error: "+msg)
	if err != nil {
		return resp
	}

	resp.Answers = append(resp.Answers, *rr)
	return resp
}

func Route(req *dns.Message) *dns.Message {
	if len(req.Questions) == 0 {
		return errorResponse(req, "no question")
	}

	q := req.Questions[0]

	if q.Type != dns.TypeTXT {
		return errorResponse(req, "only TXT queries supported")
	}

	name := strings.TrimSuffix(q.Name, ".")
	labels := strings.SplitN(name, ".", 3)

	if len(labels) < 2 {
		return errorResponse(req, "format: <command>.<input>")
	}

	command := strings.ToLower(labels[0])
	input := strings.ReplaceAll(labels[1], "-", " ")

	return handleCMD(req, command, input)
}
