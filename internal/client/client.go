// Package client fetches activity events from the GitHub API.
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github-activity/internal/model"
)

// DefaultBaseURL is the public GitHub REST API endpoint.
const DefaultBaseURL = "https://api.github.com"

// Client talks to the GitHub API.
type Client struct {
	baseURL string
	http    *http.Client
}

// New returns a Client pointing at the public GitHub API.
func New() *Client {
	return NewWithBaseURL(DefaultBaseURL)
}

// NewWithBaseURL returns a Client pointing at a custom API base URL.
func NewWithBaseURL(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 15 * time.Second},
	}
}

// FetchUserEvents returns the recent public activity for username.
func (c *Client) FetchUserEvents(ctx context.Context, username string) ([]model.Event, error) {
	endpoint := fmt.Sprintf("%s/users/%s/events", c.baseURL, url.PathEscape(username))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "github-activity-cli")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("requesting events: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apiError(resp)
	}

	var events []model.Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return events, nil
}

// apiError turns a non-200 response into a helpful error message.
func apiError(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusNotFound:
		return fmt.Errorf("user not found (HTTP %d)", resp.StatusCode)
	case http.StatusForbidden:
		if resp.Header.Get("X-RateLimit-Remaining") == "0" {
			return fmt.Errorf("GitHub API rate limit exceeded, try again later")
		}

		return fmt.Errorf("access forbidden (HTTP %d)", resp.StatusCode)
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			return fmt.Errorf("unexpected response (HTTP %d)", resp.StatusCode)
		}

		return fmt.Errorf("unexpected response (HTTP %d): %s", resp.StatusCode, msg)
	}
}
