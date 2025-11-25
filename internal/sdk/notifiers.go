package sdk

import (
	"errors"
	"fmt"

	"github.com/thecoretg/ticketbot/internal/models"
)

func (c *Client) ListNotifiers() ([]models.Notifier, error) {
	return GetMany[models.Notifier](c, "notifiers", nil)
}

func (c *Client) GetNotifier(id int) (*models.Notifier, error) {
	if id == 0 {
		return nil, errors.New("no id provided")
	}

	return GetOne[models.Notifier](c, fmt.Sprintf("notifiers/%d", id), nil)
}

func (c *Client) CreateNotifier(payload *models.Notifier) (*models.Notifier, error) {
	n := &models.Notifier{}
	if err := c.Post("notifiers", payload, n); err != nil {
		return nil, fmt.Errorf("posting to server: %w", err)
	}

	return n, nil
}

func (c *Client) DeleteNotifier(id int) error {
	if id == 0 {
		return errors.New("no id provided")
	}

	return c.Delete(fmt.Sprintf("notifiers/%d", id))
}
