package ai

import (
	"api/pkg/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AnthropicOptions struct {
	APIKey  string          `json:"-"`
	APIURL  string          `json:"api_url"`
	Model   string          `json:"model"`
	Timeout config.Duration `json:"timeout"`
}

type anthropicClient struct {
	opts       AnthropicOptions
	httpClient *http.Client
}

func newAnthropicClient(opts *AnthropicOptions) Client {
	return &anthropicClient{
		opts: *opts,
		httpClient: &http.Client{
			Timeout: opts.Timeout.Duration,
		},
	}
}

func (c *anthropicClient) IsConfigured() bool {
	return c.opts.APIKey != ""
}

func (c *anthropicClient) Complete(ctx context.Context, prompt string) (string, error) {
	if c.opts.APIKey == "" {
		return "", fmt.Errorf("AI client not configured: missing API key")
	}

	reqBody := request{
		Model:     c.opts.Model,
		MaxTokens: 2048,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.opts.APIURL, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.opts.APIKey)
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
