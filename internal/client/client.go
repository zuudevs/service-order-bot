/**

 filename  : client.go
 author    : zuudevs (zuudevs@gmail.com)
 version   : 0.1.0
 date      : 2026-05-29

 brief     : HTTP client to communicate with service-order-api

 copyright Copyright (c) 2026

**/

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIClient wraps HTTP communication with the service-order-api
type APIClient struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// New creates a new APIClient
func New(baseURL, token string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// do performs an authenticated HTTP request and decodes the JSON response
func (c *APIClient) do(method, path string, body any, out any) error {
	var bodyReader io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// GET performs a GET request
func (c *APIClient) GET(path string, out any) error {
	return c.do(http.MethodGet, path, nil, out)
}

// POST performs a POST request
func (c *APIClient) POST(path string, body any, out any) error {
	return c.do(http.MethodPost, path, body, out)
}

// PATCH performs a PATCH request
func (c *APIClient) PATCH(path string, body any, out any) error {
	return c.do(http.MethodPatch, path, body, out)
}

// DELETE performs a DELETE request
func (c *APIClient) DELETE(path string) error {
	return c.do(http.MethodDelete, path, nil, nil)
}

// Health checks the API health
func (c *APIClient) Health() error {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed: status %d", resp.StatusCode)
	}

	return nil
}