package ticket

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("ticket not found")
)

type Ticket struct {
	ID        int       `json:"id"`
	Summary   string    `json:"summary"`
	BoardID   int       `json:"board_id"`
	OwnerID   *int      `json:"owner_id"`
	CompanyID int       `json:"company_id"`
	ContactID *int      `json:"contact_id"`
	Resources *string   `json:"resources"`
	UpdatedBy *string   `json:"updated_by"`
	UpdatedOn time.Time `json:"updated_on"`
	AddedOn   time.Time `json:"added_on"`
}

type Repository interface {
	List(ctx context.Context) ([]Ticket, error)
	Get(ctx context.Context, id int) (Ticket, error)
	Upsert(ctx context.Context, c Ticket) (Ticket, error)
	Delete(ctx context.Context, id int) error
}
