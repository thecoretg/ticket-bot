package cwcompany

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("company not found")
)

type Company struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UpdatedOn time.Time `json:"updated_on"`
	AddedOn   time.Time `json:"added_on"`
}

type Repository interface {
	List(ctx context.Context) ([]Company, error)
	Get(ctx context.Context, id int) (Company, error)
	Upsert(ctx context.Context, c Company) (Company, error)
	Delete(ctx context.Context, id int) error
}
