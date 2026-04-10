package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client interface {
	// Complete sends a prompt and returns the text response.
	Complete(ctx context.Context, prompt string) (string, error)
	// IsConfigured returns whether the client has a valid API key.
	IsConfigured() bool
}

type anthropicClient struct {
	config     Config
	httpClient *http.Client
}

func NewClient(cfg *Config) Client {
	return &anthropicClient{
		config: *cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (c *anthropicClient) IsConfigured() bool {
	return c.config.APIKey != ""
}

func (c *anthropicClient) Complete(ctx context.Context, prompt string) (string, error) {
	if c.config.APIKey == "" {
		return "", fmt.Errorf("AI client not configured: missing API key")
	}

	reqBody := request{
		Model:     c.config.Model,
		MaxTokens: 2048,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.APIURL, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("AI request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp response
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse AI response: %w", err)
	}

	if len(apiResp.Content) == 0 {
		return "", fmt.Errorf("empty response from AI")
	}

	return apiResp.Content[0].Text, nil
}

// API types

type request struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type response struct {
	Content []contentBlock `json:"content"`
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
