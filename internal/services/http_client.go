package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Logger interface for logging HTTP operations
type Logger interface {
	Printf(format string, args ...interface{})
	Println(args ...interface{})
}

// HTTPClient provides HTTP client functionality for external API calls
type HTTPClient struct {
	client  *http.Client
	logger  Logger
	baseURL string
}

// HTTPResponse represents an HTTP response
type HTTPResponse struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// NewHTTPClient creates a new HTTP client with timeout and logging
func NewHTTPClient(timeout time.Duration, logger Logger) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

// Get performs a GET request to the specified URL
func (h *HTTPClient) Get(ctx context.Context, url string, headers map[string]string) (*HTTPResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("so'rov yaratishda xatolik: %w", err)
	}

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add User-Agent for better API compatibility
	req.Header.Set("User-Agent", "YordamchiDevBot/1.0")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("so'rov yuborishda xatolik: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("javobni o'qishda xatolik: %w", err)
	}

	h.logger.Printf("🌐 HTTP GET %s - Status: %d, Size: %d bytes", 
		url, resp.StatusCode, len(body))

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    resp.Header,
	}, nil
}

// GetJSON performs a GET request and unmarshals JSON response
func (h *HTTPClient) GetJSON(ctx context.Context, url string, headers map[string]string, target interface{}) error {
	resp, err := h.Get(ctx, url, headers)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP xatolik: %d", resp.StatusCode)
	}

	if err := json.Unmarshal(resp.Body, target); err != nil {
		return fmt.Errorf("JSON ni parsing qilishda xatolik: %w", err)
	}

	return nil
}