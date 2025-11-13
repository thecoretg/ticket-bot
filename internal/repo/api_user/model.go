package apiuser

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("api user not found")
)

type APIUser struct {
	ID           int       `json:"id"`
	EmailAddress string    `json:"email_address"`
	CreatedOn    time.Time `json:"created_on"`
	UpdatedOn    time.Time `json:"updated_on"`
}

type Repository interface {
	List(ctx context.Context) ([]APIUser, error)
	Get(ctx context.Context, id int) (APIUser, error)
	GetByEmail(ctx context.Context, email string) (APIUser, error)
	Insert(ctx context.Context, email string) (APIUser, error)
	Delete(ctx context.Context, id int) error
}
