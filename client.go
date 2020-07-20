package splunk

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

const (
	authContextSuffix = "/services/authentication/current-context"
)

// Client is a splunk API client
type Client struct {
	authHeader string

	config *Config
}

// Config for a SNOW connection
//
// If HTTPClient is nil, the default will be used
type Config struct {
	HTTPClient *http.Client
	BaseURL    string
}

// NewClient Creates and new splunk api client
func NewClient(ctx context.Context, username, password string, config *Config) (*Client, error) {
	configCopy := *config
	c := &Client{
		config: &configCopy,
	}
	if c.config.HTTPClient == nil {
		c.config.HTTPClient = http.DefaultClient
	}

	// Convert username:password to auth header
	c.authHeader = base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", username, password)),
	)

	// Perform simple request to make sure login is valid
	resp, err := c.MakeURLRequest(ctx, "GET", authContextSuffix, nil)
	if err != nil {
		return nil, fmt.Errorf("error making login request: %s", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad login, response code: %d", resp.StatusCode)
	}

	return c, nil
}

// MakeURLRequest Makes a request using the provided suffix and body
func (c *Client) MakeURLRequest(ctx context.Context, method, suffix string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", c.config.BaseURL, suffix), body)
	if err != nil {
		return nil, err
	}

	return c.MakeRequest(req)
}

// MakeRequest adds authentication to the request and performs it
func (c *Client) MakeRequest(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", c.authHeader))

	return c.config.HTTPClient.Do(req)
}
