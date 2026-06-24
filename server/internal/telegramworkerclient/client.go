package telegramworkerclient

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
	baseURL      string
	workerAPIKey string
	hc           *http.Client
}

type Config struct {
	BaseURL      string
	WorkerAPIKey string
}

func New(cfg Config) *Client {
	return &Client{
		baseURL:      strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/"),
		workerAPIKey: strings.TrimSpace(cfg.WorkerAPIKey),
		hc:           &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *Client) Enabled() bool {
	return c != nil && c.baseURL != ""
}

type SendRequest struct {
	Message    string `json:"message"`
	ExternalID string `json:"externalId"`
}

func (c *Client) Send(ctx context.Context, req SendRequest) error {
	if !c.Enabled() {
		return fmt.Errorf("telegram worker not configured")
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/send", bytes.NewReader(body))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.workerAPIKey != "" {
		httpReq.Header.Set("X-Channel-Worker-Key", c.workerAPIKey)
		httpReq.Header.Set("Authorization", "Bearer "+c.workerAPIKey)
	}
	res, err := c.hc.Do(httpReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("telegram worker http %d: %s", res.StatusCode, strings.TrimSpace(string(b)))
	}
	return nil
}
