package sdk

import (
	"fmt"

	"github.com/thecoretg/ticketbot/internal/server"
)

func (c *Client) GetConfig() (*server.AppConfig, error) {
	return GetOne[server.AppConfig](c, "config", nil)
}

func (c *Client) UpdateConfig(params *server.AppConfigPayload) (*server.AppConfig, error) {
	cfg := &server.AppConfig{}
	if err := c.Put("config", params, cfg); err != nil {
		return nil, fmt.Errorf("sending update request: %w", err)
	}

	return cfg, nil
}

func (c *Client) GetAppState() (*server.AppState, error) {
	return GetOne[server.AppState](c, "state", nil)
}
