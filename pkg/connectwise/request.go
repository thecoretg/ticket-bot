package connectwise

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	baseUrl = "https://api-na.myconnectwise.net/v4_6_release/apis/3.0"
)

type doRequestResult struct {
	data       []byte
	header     http.Header
	statusCode int
}

func ApiRequestNonPaginated[T any](ctx context.Context, client *Client, method, endpoint string, params *QueryParams, payload any) (*T, error) {
	if params != nil {
		endpoint = addQueryParams(endpoint, *params)
	}

	body, err := marshalBody(payload)
	if err != nil {
		return nil, err
	}

	result, err := client.doRequest(ctx, method, createFullUrl(endpoint), body)
	if err != nil {
		return nil, err
	}

	if result.statusCode < 200 || result.statusCode >= 300 {
		return nil, fmt.Errorf("bad status: %d", result.statusCode)
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
		result, err := c.doRequest(ctx, method, fullUrl, body)
		if err != nil {
			return err
		}

		if result.statusCode < 200 || result.statusCode >= 300 {
			return fmt.Errorf("bad status: %d", result.statusCode)
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
