package http

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"langgraph-sdk/schema"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

// HttpClient handles async requests to the LangGraph API.
// Adds additional error messaging & content handling above the
// provided resty client.
type HttpClient struct {
	client *resty.Client
}

// NewHttpClient creates a new HttpClient with resty.Client
func NewHttpClient(baseURL string, headers map[string]string, timeOut time.Duration, transport http.RoundTripper) *HttpClient {
	client := resty.New().
		SetBaseURL(baseURL).
		SetHeader("Accept", "application/json").
		SetHeaders(headers).
		SetTimeout(timeOut).
		SetTransport(transport)
	return &HttpClient{
		client: client,
	}
}

func (c *HttpClient) CheckConnection() error {
	_, err := c.client.R().Get("/")
	return err
}

// Get sends a GET request.
func (c *HttpClient) Get(ctx context.Context, path string, params url.Values) (*resty.Response, error) {
	req := c.client.R().SetContext(ctx)
	if params != nil {
		req.SetQueryParamsFromValues(params)
	}
	resp, err := req.Get(path)
	if err := handleError(resp, err); err != nil {
		return nil, err
	}

	return resp, nil
}

// Post sends a POST request.
func (c *HttpClient) Post(ctx context.Context, path string, jsonData any) (*resty.Response, error) {
	req := c.client.R().SetContext(ctx)

	if jsonData != nil {
		req.SetHeader("Content-Type", "application/json")
		req.SetBody(jsonData)
	}

	resp, err := req.Post(path)
	if err := handleError(resp, err); err != nil {
		return nil, err
	}

	return resp, nil
}

// Put sends a PUT request.
func (c *HttpClient) Put(ctx context.Context, path string, jsonData any) (*resty.Response, error) {
	req := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData)

	resp, err := req.Put(path)
	if err := handleError(resp, err); err != nil {
		return nil, err
	}

	return resp, nil
}

// Patch sends a PATCH request.
func (c *HttpClient) Patch(ctx context.Context, path string, jsonData any) (*resty.Response, error) {
	req := c.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData)

	resp, err := req.Patch(path)
	if err := handleError(resp, err); err != nil {
		return nil, err
	}

	return resp, nil
}

// Delete sends a DELETE request.
func (c *HttpClient) Delete(ctx context.Context, path string, jsonData any) error {
	req := c.client.R().SetContext(ctx)

	if jsonData != nil {
		req.SetHeader("Content-Type", "application/json")
		req.SetBody(jsonData)
	}

	resp, err := req.Delete(path)
	if err := handleError(resp, err); err != nil {
		return err
	}

	return nil
}

// Stream streams results using SSE.
func (c *HttpClient) Stream(ctx context.Context, path string, method string, jsonData any, params url.Values) (chan schema.StreamPart, chan error, error) {
	req := c.client.R().
		SetContext(ctx).
		SetDoNotParseResponse(true). // Important for streaming
		SetHeader("Accept", "text/event-stream").
		SetHeader("Cache-Control", "no-store")

	if jsonData != nil {
		req.SetHeader("Content-Type", "application/json")
		req.SetBody(jsonData)
	}

	if params != nil {
		req.SetQueryParamsFromValues(params)
	}

	var resp *resty.Response
	var err error

	// Execute request based on method
	switch strings.ToUpper(method) {
	case "GET":
		resp, err = req.Get(path)
	case "POST":
		resp, err = req.Post(path)
	case "PUT":
		resp, err = req.Put(path)
	case "PATCH":
		resp, err = req.Patch(path)
	case "DELETE":
		resp, err = req.Delete(path)
	default:
		return nil, nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return nil, nil, err
	}

	// Get raw response body
	rawBody := resp.RawBody()

	// Check status code
	if resp.StatusCode() >= 400 {
		// Read error body
		body, _ := io.ReadAll(rawBody)
		rawBody.Close()
		log.Printf("Error from langgraph-api: %s", string(body))
		return nil, nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode(), string(body))
	}

	// Check content type
	contentType := resp.Header().Get("Content-Type")
	if contentType == "" || !containsTextEventStream(contentType) {
		rawBody.Close()
		return nil, nil, fmt.Errorf("expected Content-Type to contain 'text/event-stream', got %s", contentType)
	}

	streamPartCh := make(chan schema.StreamPart)
	errCh := make(chan error, 1)

	// Process the SSE stream in a goroutine
	go func() {
		defer rawBody.Close()
		defer close(streamPartCh)
		defer close(errCh)

		// Parse SSE manually, since you mentioned seeing valid SSE data
		scanner := bufio.NewScanner(rawBody)
		var event, data, metadata string

		for scanner.Scan() {
			line := scanner.Text()
			// Empty line marks the end of an event
			if line == "" {
				if event != "" || data != "" {
					streamPartCh <- schema.StreamPart{
						Event:    event,
						Data:     data,
						MetaData: metadata,
					}
					// Reset for next event
					event = ""
					data = ""
					metadata = ""
				}
				continue
			} else {
				event = gjson.Get(line, "event").String()
				data = gjson.Get(line, "data").Raw
				metadata = gjson.Get(line, "metadata").Raw

				if event != "" || data != "" || metadata != "" {
					streamPartCh <- schema.StreamPart{
						Event:    event,
						Data:     data,
						MetaData: metadata,
					}
				}

				// Reset for next event
				event = ""
				data = ""
				metadata = ""
			}
		}
	}()

	return streamPartCh, errCh, nil
}
