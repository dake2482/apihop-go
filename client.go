// Package apihop provides a lightweight, zero-dependency Go client for API Hop
// (https://apihop.ai), a unified gateway for LLM, image/video generation, and
// 32+ real-time data APIs (Google Search, Twitter/X, Similarweb, Airbnb, Amazon, etc.).
package apihop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// DefaultBaseURL is the canonical API Hop gateway base URL.
const DefaultBaseURL = "https://api.apihop.ai"

// Client interacts with the API Hop REST gateway.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new API Hop client. If apiKey is empty, it falls back to
// the APIHOP_API_KEY environment variable.
func NewClient(apiKey string) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("APIHOP_API_KEY")
	}
	return &Client{
		BaseURL: DefaultBaseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// APIError represents an error payload returned by API Hop.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("apihop: http %d: %s", e.StatusCode, e.Body)
}

// DoJSON sends an authenticated JSON request to any API Hop endpoint.
func (c *Client) DoJSON(ctx context.Context, method, path string, query url.Values, reqBody any, out any) error {
	if c.APIKey == "" {
		return fmt.Errorf("apihop: missing API key (set APIHOP_API_KEY or pass to NewClient)")
	}
	base := strings.TrimRight(c.BaseURL, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	fullURL := base + path
	if len(query) > 0 {
		fullURL += "?" + query.Encode()
	}

	var bodyReader io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("apihop: marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Accept", "application/json")
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	hc := c.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	resp, err := hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return &APIError{StatusCode: resp.StatusCode, Body: string(raw)}
	}
	if out != nil && len(raw) > 0 {
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("apihop: decode response: %w", err)
		}
	}
	return nil
}

// ChatMessage represents a single chat message.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest is the request payload for POST /v1/chat-completions.
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature *float64      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
}

// ListModels calls GET /v1/models and returns the decoded JSON map.
func (c *Client) ListModels(ctx context.Context) (map[string]any, error) {
	var out map[string]any
	err := c.DoJSON(ctx, http.MethodGet, "/v1/models", nil, nil, &out)
	return out, err
}

// ChatCompletions calls POST /v1/chat-completions.
func (c *Client) ChatCompletions(ctx context.Context, req ChatRequest) (map[string]any, error) {
	var out map[string]any
	err := c.DoJSON(ctx, http.MethodPost, "/v1/chat-completions", nil, req, &out)
	return out, err
}

// TwitterUserByUsername calls GET /v1/twitter/users?username=<username>.
func (c *Client) TwitterUserByUsername(ctx context.Context, username string) (map[string]any, error) {
	q := url.Values{}
	q.Set("username", strings.TrimPrefix(username, "@"))
	var out map[string]any
	err := c.DoJSON(ctx, http.MethodGet, "/v1/twitter/users", q, nil, &out)
	return out, err
}

// TwitterUserTweets calls GET /v1/twitter/users/{id}/tweets with optional pagination_token.
func (c *Client) TwitterUserTweets(ctx context.Context, userID, paginationToken string) (map[string]any, error) {
	q := url.Values{}
	if paginationToken != "" {
		q.Set("pagination_token", paginationToken)
	}
	var out map[string]any
	err := c.DoJSON(ctx, http.MethodGet, "/v1/twitter/users/"+url.PathEscape(userID)+"/tweets", q, nil, &out)
	return out, err
}

// GetJob calls GET /v1/jobs/{id} to poll asynchronous image/video generation tasks.
func (c *Client) GetJob(ctx context.Context, jobID string) (map[string]any, error) {
	var out map[string]any
	err := c.DoJSON(ctx, http.MethodGet, "/v1/jobs/"+url.PathEscape(jobID), nil, nil, &out)
	return out, err
}
