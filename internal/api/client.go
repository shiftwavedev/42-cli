package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

// Client is a configured HTTP client for 42 API requests
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new API client with default configuration
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.intra.42.fr",
	}
}

// Get performs a GET request with authentication
func (c *Client) Get(endpoint string, token string, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return sanitizeError(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return sanitizeError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return sanitizeError(fmt.Errorf("API request failed with status: %d, response: %s", resp.StatusCode, string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return sanitizeError(err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return sanitizeError(err)
	}

	return nil
}

// PostForm performs a POST request with form data
func (c *Client) PostForm(endpoint string, data url.Values, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	resp, err := c.httpClient.PostForm(url, data)
	if err != nil {
		return sanitizeError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return sanitizeError(fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return sanitizeError(err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return sanitizeError(err)
	}

	return nil
}

// sanitizeError masks sensitive data in error messages
func sanitizeError(err error) error {
	if err == nil {
		return nil
	}

	msg := err.Error()

	// Mask Bearer tokens (common in Authorization headers)
	bearerPattern := regexp.MustCompile(`Bearer\s+[A-Za-z0-9_-]+`)
	msg = bearerPattern.ReplaceAllString(msg, "Bearer ****")

	// Mask JWT tokens (eyJ prefix)
	jwtPattern := regexp.MustCompile(`eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]*`)
	msg = jwtPattern.ReplaceAllString(msg, "****")

	// Mask client secrets (s-s4t2ud-xxx format)
	secretPattern := regexp.MustCompile(`s-s4t2ud-[a-z0-9]+`)
	msg = secretPattern.ReplaceAllString(msg, "s-s4t2ud-****")

	// Mask client IDs (u-s4t2ud-xxx format)
	clientIDPattern := regexp.MustCompile(`u-s4t2ud-[a-z0-9]+`)
	msg = clientIDPattern.ReplaceAllString(msg, "u-s4t2ud-****")

	return fmt.Errorf("%s", msg)
}

// DefaultClient is the package-level client instance
var DefaultClient = NewClient()
