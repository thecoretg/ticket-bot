package addigy

import (
	"resty.dev/v3"
)

const (
	authHeader = "x-api-key"
	baseURL    = "https://api.addigy.com/api/v2"
)

type (
	Client struct {
		restClient     *resty.Client
		defaultPage    int
		defaultPerPage int
	}

	ClientParams struct {
		Token          string
		DefaultPage    int
		DefaultPerPage int
	}
)

func NewClient(p ClientParams) *Client {
	pg := p.DefaultPage
	if pg == 0 {
		pg = 1
	}

	perPg := p.DefaultPerPage
	if perPg == 0 {
		perPg = 50
	}

	c := resty.New()
	c.SetHeader(authHeader, p.Token)
	c.SetHeader("Content-Type", "application/json")
	c.SetHeader("Accept", "application/json")
	c.SetRetryCount(3)
	c.SetBaseURL(baseURL)

	return &Client{
		restClient:     c,
		defaultPage:    pg,
		defaultPerPage: perPg,
	}
}
