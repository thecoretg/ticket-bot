package sdk

import "github.com/thecoretg/ticketbot/internal/db"

func (c *Client) SyncRooms() error {
	return c.Post("sync/webex_rooms", nil, nil)
}

func (c *Client) ListRooms() ([]db.WebexRoom, error) {
	return GetMany[db.WebexRoom](c, "rooms", nil)
}
