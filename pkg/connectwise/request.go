package connectwise

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseUrl = "https://api-na.myconnectwise.net/v4_6_release/apis/3.0"
)

type doRequestResult struct {
	data       []byte
	header     http.Header
	statusCode int
}

type ErrNotFound struct{}

func (e *ErrNotFound) Error() string {
	return "not found"
}

type ErrBadRequest struct {
	Message string
}

func (e *ErrBadRequest) Error() string {
	return fmt.Sprintf("bad request: %s", e.Message)
}

type ErrMaxRetries struct {
	Attempts  int
	LastError error
}

func (e *ErrMaxRetries) Error() string {
	return fmt.Sprintf("max retries exceeded after %d attempts: %v", e.Attempts, e.LastError)
}

func ApiRequestNonPaginated[T any](ctx context.Context, client *Client, method, endpoint string, params *QueryParams, payload any) (*T, error) {
	if params != nil {
		endpoint = addQueryParams(endpoint, *params)
	}

	body, err := marshalBody(payload)
	if err != nil {
		return nil, err
	}

	result, err := client.doRequestWithRetry(ctx, method, createFullUrl(endpoint), body)
	if err != nil {
		return nil, err
	}

	if result.statusCode < 200 || result.statusCode >= 300 {
		return nil, handleErrorResponse(result)
	}

	if len(result.data) == 0 {
		var zero T
		return &zero, nil
	}

	var target T
	if err := json.Unmarshal(result.data, &target); err != nil {
		return nil, fmt.Errorf("unmarshaling the response to json: %w", err)
	}

	return &target, nil
}

func ApiRequestPaginated[T any](ctx context.Context, c *Client, method, endpoint string, params *QueryParams, payload any) ([]T, error) {
	if params != nil {
		endpoint = addQueryParams(endpoint, *params)
	}

	body, err := marshalBody(payload)
	if err != nil {
		return nil, err
	}

	var allItems []T
	err = c.doPaginatedRequest(ctx, method, createFullUrl(endpoint), body, func(data []byte) error {
		var pageItems []T
		if err := json.Unmarshal(data, &pageItems); err != nil {
			return fmt.Errorf("unmarshaling page: %w", err)
		}
		allItems = append(allItems, pageItems...)
		return nil
	})

	return allItems, err
}

func (c *Client) doPaginatedRequest(ctx context.Context, method, fullUrl string, body io.Reader, handlePage func(data []byte) error) error {
	for {
		result, err := c.doRequestWithRetry(ctx, method, fullUrl, body)
		if err != nil {
			return err
		}

		if result.statusCode < 200 || result.statusCode >= 300 {
			return handleErrorResponse(result)
		}

		if err := handlePage(result.data); err != nil {
			return err
		}

		linkHeader := result.header.Get("Link")
		nextUrl, found := parseLinkHeader(linkHeader, "next")
		if !found {
			break
		}
		fullUrl = nextUrl
	}

	return nil
}

func (c *Client) doRequestWithRetry(ctx context.Context, method, fullUrl string, body io.Reader) (*doRequestResult, error) {
	var lastErr error
	var result *doRequestResult

	for attempt := 0; attempt < c.retryConfig.MaxRetries; attempt++ {
		var requestBody io.Reader
		if body != nil {
			requestBody = body
		}
		result, lastErr = c.doRequest(ctx, method, fullUrl, requestBody)

		if lastErr == nil && !isRetryableError(result.statusCode, nil) {
			return result, nil
		}

		if attempt == c.retryConfig.MaxRetries {
			break
		}

		var delay time.Duration
		if result != nil && result.statusCode == http.StatusTooManyRequests {
			if retryAfter := extractRetryAfter(result.header); retryAfter > 0 {
				delay = retryAfter
			} else {
				delay = calculateBackoff(attempt, c.retryConfig)
			}
		} else {
			delay = calculateBackoff(attempt, c.retryConfig)
			// back off if you know what's good for ya kid
		}

		select {
		case <-ctx.Done():
			slog.Debug("request cancelled", "method", method, "url", fullUrl, "attempt", attempt+1, "error", ctx.Err())
			return nil, ctx.Err()
		case <-time.After(delay):
			// continue to the next attempt
			slog.Debug("retrying request", "method", method, "url", fullUrl, "attempt", attempt+1, "delay", delay)
		}
	}

	if lastErr != nil {
		return nil, &ErrMaxRetries{
			Attempts:  c.retryConfig.MaxRetries + 1,
			LastError: lastErr,
		}
	}

	return result, nil
}

func handleErrorResponse(result *doRequestResult) error {
	switch result.statusCode {
	case http.StatusNotFound:
		return &ErrNotFound{}
	case http.StatusBadRequest:
		var errMsg string
		if len(result.data) > 0 {
			errMsg = string(result.data)
		}
		return &ErrBadRequest{Message: errMsg}
	default:
		return fmt.Errorf("unexpected status code: %d, response: %s", result.statusCode, string(result.data))
	}
}

func (c *Client) doRequest(ctx context.Context, method, fullUrl string, body io.Reader) (*doRequestResult, error) {
	req, err := http.NewRequestWithContext(ctx, method, fullUrl, body)
	if err != nil {
		return nil, fmt.Errorf("creating the request: %w", err)
	}

	c.setStandardHeaders(req)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending the request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// TODO: handle this
		}
	}(res.Body)

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	return &doRequestResult{
		data:       data,
		header:     res.Header,
		statusCode: res.StatusCode,
	}, nil
}

func (c *Client) setStandardHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("clientId", c.clientId)
	req.Header.Set("Authorization", c.encodedCreds)
}

func marshalBody(body any) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling the body: %w", err)
	}

	return strings.NewReader(string(data)), nil
}

func parseLinkHeader(linkHeader, rel string) (string, bool) {
	links := strings.Split(linkHeader, ",")
	for _, link := range links {
		parts := strings.Split(strings.TrimSpace(link), ";")
		if len(parts) < 2 {
			continue
		}
		urlPart := strings.Trim(parts[0], "<>")
		relPart := strings.TrimSpace(parts[1])
		if relPart == fmt.Sprintf(`rel="%s"`, rel) {
			return urlPart, true
		}
	}

	return "", false
}

func addQueryParams(endpoint string, params QueryParams) string {
	u, _ := url.Parse(endpoint)
	p := params.ToMap()
	q := u.Query()
	for k, v := range p {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func createFullUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s", baseUrl, endpoint)
}

func isRetryableError(statusCode int, err error) bool {
	if err != nil {
		return true
	}

	switch statusCode {
	case http.StatusTooManyRequests,
		http.StatusRequestTimeout,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func calculateBackoff(attempt int, config *RetryConfig) time.Duration {
	delay := float64(config.InitialDelay) * math.Pow(config.BackOffMultiplier, float64(attempt))
	if delay > float64(config.MaxDelay) {
		delay = float64(config.MaxDelay)
	}
	return time.Duration(delay)
}

func extractRetryAfter(headers http.Header) time.Duration {
	retryAfter := headers.Get("Retry-After")
	if retryAfter == "" {
		return 0
	}

	if seconds, err := time.ParseDuration(retryAfter + "s"); err == nil {
		return seconds
	}

	return 0
}
