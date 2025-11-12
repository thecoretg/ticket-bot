package contact

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("contact not found")
)

type Contact struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  *string   `json:"last_name"`
	CompanyID *int      `json:"company_id"`
	UpdatedOn time.Time `json:"updated_on"`
	AddedOn   time.Time `json:"added_on"`
}

type Repository interface {
	List(ctx context.Context) ([]Contact, error)
	Get(ctx context.Context, id int) (Contact, error)
	Upsert(ctx context.Context, c Contact) (Contact, error)
	Delete(ctx context.Context, id int) error
}
