package handlers

import (
	"dnsense/ai"
	"dnsense/dns"
	"fmt"
)

type command struct {
	prompt string
}

var commands = map[string]command{
	"explain": {
		prompt: "Explain '%s' in 2-3 sentences. Be concise, technical, and clear. No markdown.",
	},
	"tldr": {
		prompt: "Give a TL;DR of '%s' in one sentence, max 200 characters. No markdown.",
	},
	"define": {
		prompt: "Define '%s' in one sentence like a dictionary. Max 200 characters. No markdown.",
	},
	"ask": {
		prompt: "Answer this question factually and concisely in 2-3 sentences. Be precise. No markdown. No preamble. Question: '%s'",
	},
}

func handleCMD(req *dns.Message, cmdName string, input string) *dns.Message {
	cmd, ok := commands[cmdName]
	if !ok {
		return errorResponse(req, "unknown command: "+cmdName)
	}

	text, err := ai.Ask(fmt.Sprintf(cmd.prompt, input))
	if err != nil {
		return errorResponse(req, "ai error: "+err.Error())
	}

	rr, err := dns.NewTXTRecord(req.Questions[0].Name, 300, text)
	if err != nil {
		return errorResponse(req, "encoding error")
	}

	resp := baseResponse(req)
	resp.Answers = append(resp.Answers, *rr)
	return resp
}
