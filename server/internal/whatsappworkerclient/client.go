package whatsappworkerclient

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
	Message       string `json:"message,omitempty"`
	ExternalID    string `json:"externalId,omitempty"`
	DestinationID string `json:"destinationId,omitempty"`
	Type          string `json:"type,omitempty"`
	VideoID       string `json:"videoId,omitempty"`
	TemplateSlug  string `json:"templateSlug,omitempty"`
}

type StatusResponse struct {
	Connected bool `json:"connected"`
}

type QRResponse struct {
	QR string `json:"qr"`
}

func (c *Client) Send(ctx context.Context, req SendRequest) error {
	if !c.Enabled() {
		return fmt.Errorf("whatsapp worker not configured")
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
		httpReq.Header.Set("X-Worker-Key", c.workerAPIKey)
		httpReq.Header.Set("Authorization", "Bearer "+c.workerAPIKey)
	}
	res, err := c.hc.Do(httpReq)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("whatsapp worker http %d: %s", res.StatusCode, strings.TrimSpace(string(b)))
	}
	return nil
}

func (c *Client) Status(ctx context.Context) (*StatusResponse, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("whatsapp worker not configured")
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/status", nil)
	if err != nil {
		return nil, err
	}
	res, err := c.hc.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("whatsapp worker status http %d", res.StatusCode)
	}
	var out StatusResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) QR(ctx context.Context) (*QRResponse, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("whatsapp worker not configured")
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/qr", nil)
	if err != nil {
		return nil, err
	}
	res, err := c.hc.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("whatsapp worker qr http %d", res.StatusCode)
	}
	var out QRResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
