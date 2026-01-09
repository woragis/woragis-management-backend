package aiservice

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client is an HTTP client for the AI Service
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new AI Service client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // AI requests can take longer
		},
	}
}

// ChatRequest represents a chat request to the AI service
type ChatRequest struct {
	Agent      string   `json:"agent"`
	Input      string   `json:"input"`
	System     *string  `json:"system,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	Model      *string  `json:"model,omitempty"`
	Provider   *string  `json:"provider,omitempty"`
}

// ChatResponse represents a chat response from the AI service
type ChatResponse struct {
	Agent  string `json:"agent"`
	Output string `json:"output"`
}

// ChatStreamRequest represents a streaming chat request
type ChatStreamRequest struct {
	ChatRequest
}

// Chat sends a chat request to the AI service
func (c *Client) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	url := fmt.Sprintf("%s/v1/chat", c.baseURL)

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	var response ChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// ChatStream sends a streaming chat request to the AI service
// The onChunk callback is called for each chunk received
func (c *Client) ChatStream(ctx context.Context, req ChatStreamRequest, onChunk func(string)) error {
	url := fmt.Sprintf("%s/v1/chat/stream", c.baseURL)

	jsonData, err := json.Marshal(req.ChatRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read streaming response (NDJSON format: one JSON object per line)
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		var chunk struct {
			Delta string `json:"delta,omitempty"`
			Done  bool   `json:"done,omitempty"`
			Error string `json:"error,omitempty"`
		}
		
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue // Skip invalid JSON lines
		}
		
		if chunk.Error != "" {
			return fmt.Errorf("AI service stream error: %s", chunk.Error)
		}
		
		if chunk.Done {
			break
		}
		
		if chunk.Delta != "" && onChunk != nil {
			onChunk(chunk.Delta)
		}
	}

	return scanner.Err()
}

// HealthCheck checks if the AI service is healthy
func (c *Client) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/healthz", c.baseURL)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("AI service health check failed with status %d", resp.StatusCode)
	}

	return nil
}

