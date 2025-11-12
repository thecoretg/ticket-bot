package member

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("member not found")
)

type Member struct {
	ID           int       `json:"id"`
	Identifier   string    `json:"identifier"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	PrimaryEmail string    `json:"primary_email"`
	UpdatedOn    time.Time `json:"updated_on"`
	AddedOn      time.Time `json:"added_on"`
}

type Repository interface {
	List(ctx context.Context) ([]Member, error)
	Get(ctx context.Context, id int) (Member, error)
	Upsert(ctx context.Context, c Member) (Member, error)
	Delete(ctx context.Context, id int) error
}
