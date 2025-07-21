package addigy

import "resty.dev/v3"

func NewClient(token string) *resty.Client {
	c := resty.New()
	c.SetHeaderAuthorizationKey("x-api-key")
	c.SetAuthToken(token)
	c.SetHeader("Accept", "application/json")

	return c
}
