package sdk

import "github.com/thecoretg/ticketbot/internal/db"

func (c *Client) SyncRooms() error {
	return Post(c, "sync/webex_rooms", nil)
}

func (c *Client) ListRooms() ([]db.WebexRoom, error) {
	return GetMany[db.WebexRoom](c, "rooms", nil)
}
