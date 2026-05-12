package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func Ask(prompt string) (string, error) {
	payload := map[string]any{
		"model":  "llama3-chatqa:8b",
		"prompt": prompt,
		"stream": false,
		"think":  false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	resp, err := client.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	var result struct {
		Response string `json:"response"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("unmarshal: %w", err)
	}

	return strings.TrimSpace(result.Response), nil
}
