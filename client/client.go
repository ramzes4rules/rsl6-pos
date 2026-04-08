package rslpos

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ramzes4rules/rsl6.pos.client/internal/soap"
)

// Client implements the Service interface as a SOAP 1.2 client
// for the RSLoyaltyService WCF service.
type Client struct {
	url        string
	httpClient *http.Client
}

// Ensure Client implements Service at compile time.
var _ Service = (*Client)(nil)

// Option configures the Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client (e.g. for TLS configuration or timeouts).
func WithHTTPClient(c *http.Client) Option {
	return func(client *Client) {
		client.httpClient = c
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) Option {
	return func(client *Client) {
		client.httpClient.Timeout = d
	}
}

// NewClient creates a new RSLoyaltyService SOAP client for the given URL.
func NewClient(url string, opts ...Option) *Client {
	c := &Client{
		url: url,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// doSOAP performs a SOAP 1.2 call and returns the raw response bytes.
func (c *Client) doSOAP(ctx context.Context, operationName string, requestBody interface{}) ([]byte, error) {
	action := soap.ActionPrefix + operationName
	envelope, err := soap.BuildEnvelope(action, c.url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("build soap envelope: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(envelope))
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}
	req.Header.Set("Content-Type", fmt.Sprintf(`application/soap+xml; charset=utf-8; action="%s"`, action))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusInternalServerError {
		return nil, fmt.Errorf("unexpected http status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// call performs a full SOAP call: send request, parse envelope, check faults,
// and optionally unmarshal the response body into result.
func (c *Client) call(ctx context.Context, operationName string, requestBody interface{}, result interface{}) error {
	respBody, err := c.doSOAP(ctx, operationName, requestBody)
	if err != nil {
		return err
	}

	parsed, err := soap.ParseEnvelope(respBody)
	if err != nil {
		return err
	}
	if parsed.Fault != nil {
		return &FaultError{
			Code:   parsed.Fault.Code,
			Reason: parsed.Fault.Reason,
			Detail: parsed.Fault.Detail,
		}
	}
	if result != nil {
		if err := xml.Unmarshal(parsed.BodyContent, result); err != nil {
			return fmt.Errorf("unmarshal soap body: %w", err)
		}
	}
	return nil
}
