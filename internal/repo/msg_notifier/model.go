package notifier

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("notifier not found")
)

type Notifier struct {
	ID            int       `json:"id"`
	CwBoardID     int       `json:"cw_board_id"`
	WebexRoomID   int       `json:"webex_room_id"`
	NotifyEnabled bool      `json:"notify_enabled"`
	CreatedOn     time.Time `json:"created_on"`
}

type Repository interface {
	ListAll(ctx context.Context) ([]Notifier, error)
	ListByBoard(ctx context.Context, boardID int) ([]Notifier, error)
	ListByRoom(ctx context.Context, roomID int) ([]Notifier, error)
	Get(ctx context.Context, id int) (Notifier, error)
	Insert(ctx context.Context, n Notifier) (Notifier, error)
	Update(ctx context.Context, n Notifier) (Notifier, error)
	Delete(ctx context.Context, id int) error
}
