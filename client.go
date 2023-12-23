package paystack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

const defaultBaseURL = "https://api.paystack.co"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client interface {
	ResolveBankAccount(ctx context.Context, bankCode, accountNumber string) (data ResolveResponse, err error)
}

var _ Client = (*client)(nil)

type client struct {
	httpClient HttpClient
	logger     *logrus.Logger
	baseURL    string
	token      string
}

// ClientOption is a function that configures a Client.
type ClientOption func(*client)

// WithHTTPClient sets the HTTP client for the paystack API client.
func WithHTTPClient(c HttpClient) ClientOption {
	return func(target *client) {
		target.httpClient = c
	}
}

// WithLogger sets the *logrus.Logger for the paystack API client.
func WithLogger(l *logrus.Logger) ClientOption {
	return func(target *client) {
		target.logger = l
	}
}

// WithBaseURL sets the base URL for the paystack API client.
func WithBaseURL(baseURL string) ClientOption {
	return func(target *client) {
		target.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}

func NewClient(token string, options ...ClientOption) *client {
	c := &client{
		httpClient: http.DefaultClient,
		baseURL:    defaultBaseURL,
		token:      token,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *client) newRequest(ctx context.Context, method, url string) (*request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	return NewRequest(req), nil
}

func (c *client) do(ctx context.Context, req *request) error {
	if c.logger != nil {
		c.logger.WithContext(ctx).WithFields(logrus.Fields{
			"http.request.method": req.req.Method,
			"http.request.url":    req.req.URL.String(),
		}).Debug("paystack.client -> request")
	}

	resp, err := c.httpClient.Do(req.req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(b))

	if c.logger != nil {
		c.logger.WithContext(ctx).WithFields(logrus.Fields{
			"http.response.status_code":  resp.StatusCode,
			"http.response.body.content": string(b),
			"http.response.headers":      resp.Header,
		}).Debug("paystack.client -> response")
	}

	if req.decodeTo != nil {
		if err := json.NewDecoder(resp.Body).Decode(req.decodeTo); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
