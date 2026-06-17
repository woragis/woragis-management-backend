package creativesclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	baseURL string
	apiKey  string
	hc      *http.Client
}

type Config struct {
	BaseURL string
	APIKey  string
}

func New(cfg Config) *Client {
	return &Client{
		baseURL: strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/"),
		apiKey:  strings.TrimSpace(cfg.APIKey),
		hc:      &http.Client{Timeout: 120 * time.Second},
	}
}

func (c *Client) Enabled() bool {
	return c != nil && c.baseURL != "" && c.apiKey != ""
}

type ContentImageJob struct {
	ID           uuid.UUID `json:"id"`
	ExternalID   *string   `json:"externalId"`
	Status       string    `json:"status"`
	OutputURL    *string   `json:"outputUrl"`
	ErrorMessage *string   `json:"errorMessage"`
}

type CreateContentImageRequest struct {
	ExternalID    string
	Mode          string
	Model         string
	Prompt        string
	Size          string
	Quality       string
	ReferenceURLs []string
	WebhookURL    *string
	Metadata      map[string]any
}

func (c *Client) CreateContentImage(ctx context.Context, req CreateContentImageRequest) (*ContentImageJob, error) {
	body, err := json.Marshal(map[string]any{
		"externalId":    req.ExternalID,
		"mode":          req.Mode,
		"model":         req.Model,
		"prompt":        req.Prompt,
		"size":          req.Size,
		"quality":       req.Quality,
		"referenceUrls": req.ReferenceURLs,
		"webhookUrl":    req.WebhookURL,
		"metadata":      req.Metadata,
	})
	if err != nil {
		return nil, err
	}
	var out ContentImageJob
	if _, err := c.doJSON(ctx, http.MethodPost, "/v1/content-images", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetContentImage(ctx context.Context, id uuid.UUID) (*ContentImageJob, error) {
	var out ContentImageJob
	_, err := c.doJSON(ctx, http.MethodGet, "/v1/content-images/"+id.String(), nil, &out)
	return &out, err
}

func (c *Client) DownloadOutput(ctx context.Context, outputURL string) ([]byte, error) {
	url := outputURL
	if strings.HasPrefix(outputURL, "/") {
		url = c.baseURL + outputURL
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		raw, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("download output http %d: %s", res.StatusCode, string(raw))
	}
	return io.ReadAll(res.Body)
}

func (c *Client) doJSON(ctx context.Context, method, path string, body []byte, out any) (*http.Response, error) {
	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if out == nil {
			_ = res.Body.Close()
		}
	}()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return res, fmt.Errorf("creatives http %d: %s", res.StatusCode, string(raw))
	}
	if out != nil {
		if err := json.Unmarshal(raw, out); err != nil {
			return res, err
		}
	}
	return res, nil
}
