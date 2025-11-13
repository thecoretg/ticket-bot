package userfwd

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("forward rule not found")
)

type Forward struct {
	ID            int        `json:"id"`
	UserEmail     string     `json:"user_email"`
	DestRoomID    int        `json:"dest_room_id"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	Enabled       bool       `json:"enabled"`
	UserKeepsCopy bool       `json:"user_keeps_copy"`
	UpdatedOn     time.Time  `json:"updated_on"`
	AddedOn       time.Time  `json:"added_on"`
}

type Repository interface {
	ListAll(ctx context.Context) ([]Forward, error)
	ListByEmail(ctx context.Context, email string) ([]Forward, error)
	Get(ctx context.Context, id int) (Forward, error)
	Insert(ctx context.Context, c Forward) (Forward, error)
	Delete(ctx context.Context, id int) error
}
