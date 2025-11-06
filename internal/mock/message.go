package mock

import (
	webex2 "github.com/thecoretg/ticketbot/internal/external/webex"
)

// This isn't ideal at all, but I need it short term just to get this up and running.
// When I get time I will make a mock setup that doesn't actually call the Webex API.

type WebexClient struct {
	// real webex client used for listing rooms, but not posting messages
	webexClient *webex2.Client
}

func NewWebexClient(token string) *WebexClient {
	return &WebexClient{
		webexClient: webex2.NewClient(token),
	}
}

func (w *WebexClient) PostMessage(message *webex2.Message) (*webex2.Message, error) {
	return message, nil
}

func (w *WebexClient) ListRooms(params map[string]string) ([]webex2.Room, error) {
	return w.webexClient.ListRooms(params)
}
