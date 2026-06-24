package agentworkerclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	agentAPIKey string
	hc         *http.Client
}

type Config struct {
	BaseURL     string
	AgentAPIKey string
}

func New(cfg Config) *Client {
	return &Client{
		baseURL:     strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/"),
		agentAPIKey: strings.TrimSpace(cfg.AgentAPIKey),
		hc:          &http.Client{Timeout: 90 * time.Second},
	}
}

func (c *Client) Enabled() bool {
	return c != nil && c.baseURL != ""
}

type ComposeRequest struct {
	TemplateBody       string         `json:"templateBody"`
	ComposeMode        string         `json:"composeMode,omitempty"`
	Data               map[string]any `json:"data,omitempty"`
	DestinationContext map[string]any `json:"destinationContext,omitempty"`
}

type ComposeResponse struct {
	Message string `json:"message"`
}

func (c *Client) Compose(ctx context.Context, req ComposeRequest) (string, error) {
	if !c.Enabled() {
		return "", fmt.Errorf("agent worker not configured")
	}
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/automated/compose", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.agentAPIKey != "" {
		httpReq.Header.Set("X-Agent-Key", c.agentAPIKey)
		httpReq.Header.Set("Authorization", "Bearer "+c.agentAPIKey)
	}
	res, err := c.hc.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("agent worker http %d: %s", res.StatusCode, strings.TrimSpace(string(b)))
	}
	var out ComposeResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.Message), nil
}
