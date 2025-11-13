package cwboard

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("board not found")
)

type Board struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	UpdatedOn time.Time `json:"updated_on"`
	AddedOn   time.Time `json:"added_on"`
}

type Repository interface {
	List(ctx context.Context) ([]Board, error)
	Get(ctx context.Context, id int) (Board, error)
	Upsert(ctx context.Context, b Board) (Board, error)
	Delete(ctx context.Context, id int) error
}
