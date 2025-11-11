package sdk

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"resty.dev/v3"
)

type Client struct {
	restClient *resty.Client
}

func NewClient(apiKey, baseURL string) (*Client, error) {
	c := resty.New()
	c.SetBaseURL(baseURL)
	c.SetHeader("Content-Type", "application/json")
	c.SetHeader("Accept", "application/json")
	c.SetRetryCount(3)

	if apiKey != "" {
		c.SetAuthToken(apiKey)
	}

	return &Client{restClient: c}, nil
}

var (
	ErrNotFound = errors.New("404 status returned")
)

func (c *Client) Ping() error {
	res, err := c.restClient.R().
		Get("")

	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf("error response from ticketbot api: %w", err)
	}

	return nil
}

func (c *Client) AuthTest() error {
	res, err := c.restClient.R().
		Get("authtest")

	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf("error response from ticketbot api: %s", res.String())
	}

	return nil
}

func GetOne[T any](c *Client, endpoint string, params map[string]string) (*T, error) {
	var target T
	res, err := c.restClient.R().
		SetQueryParams(params).
		SetResult(&target).
		Get(endpoint)

	if err != nil {
		return nil, err
	}

	if res.IsError() {
		if res.StatusCode() == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error response from ConnectWise API: %s", res.String())
	}

	return res.Result().(*T), nil
}

func GetMany[T any](c *Client, endpoint string, params map[string]string) ([]T, error) {
	var allItems []T

	for endpoint != "" {
		var target []T
		req := c.restClient.R().
			SetQueryParams(params).
			SetResult(&target)

		res, err := req.Get(endpoint)
		if err != nil {
			return nil, err
		}

		if res.IsError() {
			if res.StatusCode() == http.StatusNotFound {
				return nil, ErrNotFound
			}
			return nil, fmt.Errorf("error response from ConnectWise API: %s", res.String())
		}

		for _, item := range target {
			allItems = append(allItems, item)
		}

		params = nil
		endpoint = parseLinkHeader(res.Header().Get("Link"), "next")
	}

	return allItems, nil
}

func (c *Client) Post(endpoint string, body, target any) error {
	req := c.restClient.R().
		SetBody(body)

	if target != nil {
		req.SetResult(target)
	}

	res, err := req.Post(endpoint)
	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf("error response from API: %s", res.String())
	}

	return nil
}

func (c *Client) Put(endpoint string, body, target any) error {
	req := c.restClient.R().
		SetBody(body)

	if target != nil {
		req.SetResult(target)
	}

	res, err := req.Put(endpoint)
	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf("error response from API: %s", res.String())
	}

	return nil
}

func (c *Client) Delete(endpoint string) error {
	res, err := c.restClient.R().
		Delete(endpoint)

	if err != nil {
		return err
	}

	if res.IsError() {
		if res.StatusCode() == http.StatusNotFound {
			return ErrNotFound
		}
		return fmt.Errorf("error response from API: %s", res.String())
	}

	return nil
}

func parseLinkHeader(linkHeader, rel string) string {
	links := strings.Split(linkHeader, ",")
	for _, link := range links {
		parts := strings.Split(strings.TrimSpace(link), ";")
		if len(parts) < 2 {
			continue
		}
		urlPart := strings.Trim(parts[0], "<>")
		relPart := strings.TrimSpace(parts[1])
		if relPart == fmt.Sprintf(`rel="%s"`, rel) {
			return urlPart
		}
	}

	return ""
}
