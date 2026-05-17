#!/bin/bash
set -e

# ── setup ────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

cleanline() { printf '\033[1A\033[2K\r'; }
ok()   { echo -e "${GREEN}✓${NC} $1"; }
info() { echo -e "${BLUE}→${NC} $1"; sleep 0.5 ; cleanline; }
warn() { echo -e "${YELLOW}⚠${NC} $1"; }
fail() { echo -e "${RED}✗${NC} $1"; exit 1; }

echo ""
echo -e "${GREEN}initiating dnsense..."
echo ""
# ──────────────────────────────────────────────────────

info "checking go..."
if ! command -v go &> /dev/null; then
    fail "go not found. install from https://go.dev/dl"
fi
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED="1.25"
if [ "$(printf '%s\n' "$REQUIRED" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED" ]; then
    fail "go $REQUIRED+ required, found $GO_VERSION"
fi 

ok "go $GO_VERSION"

info "checking ollama..."
if ! command -v ollama &> /dev/null; then
    warn "ollama not found. installing..."
    curl -fsSL https://ollama.com/install.sh | sh
    ok "ollama installed"
fi

if ! curl -sf http://localhost:11434/api/tags > /dev/null 2>&1; then
    warn "ollama not running."
		info "starting ollama service..."
    ollama serve > /tmp/ollama.log 2>&1 &
    OLLAMA_PID=$!
    echo $OLLAMA_PID > /tmp/ollama.pid

    ATTEMPTS=0
    until curl -sf http://localhost:11434/api/tags > /dev/null 2>&1; do
        sleep 1
        ATTEMPTS=$((ATTEMPTS + 1))
        if [ $ATTEMPTS -ge 30 ]; then
            fail "ollama failed to start. check /tmp/ollama.log"
        fi
    done
    ok "ollama started (pid $OLLAMA_PID)"
else
    ok "ollama already running"
fi

MODEL="llama3-chatqa:8b"
info "checking model $MODEL..."
if ! ollama list 2>/dev/null | grep -q "llama3-chatqa"; then
    warn "model not found. pulling $MODEL (~5GB, once only)..."
    ollama pull $MODEL
    ok "model $MODEL ready"
else
    ok "model $MODEL is present"
fi

info "building dnsense..."
if ! go build -ldflags="-s -w" -trimpath -o dnsense .; then
    fail "build failed"
fi
ok "dnsense build successful ($(du -sh dnsense | cut -f1) binary)"

PORT=${DNS_PORT:-8053}
info "checking port $PORT..."
if lsof -i UDP:$PORT > /dev/null 2>&1; then
    fail "port $PORT already in use. set DNS_PORT=<other> to override"
fi
ok "port $PORT is free"

echo ""
echo -e "${GREEN}checks passed. launching dnsense${NC}"
echo ""

DNS_PORT=$PORT exec ./dnsense
