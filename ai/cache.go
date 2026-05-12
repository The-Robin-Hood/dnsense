package ai

import (
	"sync"
	"time"
)

type cacheEntry struct {
	response  string
	expiresAt time.Time
}

var (
	cache   = make(map[string]cacheEntry)
	cacheMu sync.RWMutex
	cacheTTL = 5 * time.Minute
)

func getCached(prompt string) (string, bool) {
	cacheMu.RLock()
	defer cacheMu.RUnlock()

	entry, ok := cache[prompt]
	if !ok || time.Now().After(entry.expiresAt) {
		return "", false
	}
	return entry.response, true
}

func setCached(prompt, response string) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cache[prompt] = cacheEntry{
		response:  response,
		expiresAt: time.Now().Add(cacheTTL),
	}
}
