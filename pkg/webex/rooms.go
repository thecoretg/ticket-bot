package webex

import (
	"context"
	"fmt"
)

func (c *Client) ListRooms(ctx context.Context) ([]Room, error) {
	r := &ListRoomsResponse{}
	if err := c.request(ctx, "GET", "rooms", nil, r); err != nil {
		return nil, fmt.Errorf("getting rooms: %w", err)
	}

	return r.Items, nil
}

func (c *Client) GetRoom(ctx context.Context, roomID string) (*Room, error) {
	r := &Room{}
	if err := c.request(ctx, "GET", fmt.Sprintf("rooms/%s", roomID), nil, r); err != nil {
		return nil, fmt.Errorf("getting room: %w", err)
	}

	return r, nil
}
