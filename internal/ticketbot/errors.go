package ticketbot

import (
	"errors"
	"fmt"
	"log/slog"
	"tctg-automation/pkg/connectwise"
)

type ErrWasDeleted struct {
	ItemType string
	ItemID   int
}

func (e ErrWasDeleted) Error() string {
	return fmt.Sprintf("%s %d was deleted by external factors", e.ItemType, e.ItemID)
}

// checks for specific errors to reduce repetitive connectwise error checking
func checkCWError(msg, itemType string, err error, itemID int) error {
	var notFoundErr *connectwise.ErrNotFound

	switch {
	case errors.As(err, &notFoundErr):
		slog.Info("item was deleted externally", "id", itemID)
		return ErrWasDeleted{
			ItemType: itemType,
			ItemID:   itemID,
		}
	default:
		return fmt.Errorf("%s: %w", msg, err)
	}
}
