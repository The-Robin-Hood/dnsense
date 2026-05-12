# dnsense

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8)
![DNS](https://img.shields.io/badge/protocol-DNS-blueviolet)
![port](https://img.shields.io/badge/port-53-critical)

AI over DNS. Ask questions with `dig`. Get answers in TXT records.
No HTTP. No browser. No dignity.

> You came here looking for a normal project. This is not that.
> This runs an LLM through a protocol designed in 1983 to look up phone books.
---

## Usage

```bash
dig TXT ask.why-does-hair-go-white @your-server
dig TXT explain.recursion @your-server
dig TXT tldr.bitcoin @your-server
dig TXT define.ephemeral @your-server
```

> yes, hyphens instead of spaces.
> DNS said no to spaces in 1983 and has not changed its mind.

---

## commands

```
ask.<question>      anything. use hyphens.
explain.<topic>     2-3 sentence technical breakdown
tldr.<topic>        one sentence. brutally short.
define.<word>       dictionary style. no fluff.
```

---

## How it works

- you send a DNS TXT query with `dig`
- server parses the domain name as `command.input`
- fires the input at local Ollama as a prompt
- chunks the response into 255-byte TXT strings (RFC 1035 limit)
- sends it back over raw UDP
- your terminal shows the answer

> the entire internet between you and the answer is 512 bytes of UDP.

---

## Requirements

- Go 1.22+
- Ollama running locally — `ollama pull llama3-chatqa:8b` else use claude or openAI 
- a terminal. you have one. you're reading this.

---

## Run it

```bash
git clone https://github.com/The-Robin-Hood/dnsense
cd dnsense
go run .

# fire away
dig TXT ask.what-is-dark-matter @127.0.0.1 -p 8053
```

---

## Why DNS

- works through firewalls that block everything else
- no TLS handshake, no headers, no cookies, no noise
- `dig` is on every machine on earth
- responses fit in 512 bytes — answers stay tight
- nobody is watching your DNS queries. probably.

---

## Built from scratch

zero DNS libraries. raw bytes. RFC 1035 followed line by line.
header parsing, name encoding, compression pointers, TXT record
chunking.

> because using `miekg/dns` would have been too easy
> and I really wanted to learn how this thing works.

---

*rule 53 — if it exists, someone has done it with DNS. this is that someone.*
