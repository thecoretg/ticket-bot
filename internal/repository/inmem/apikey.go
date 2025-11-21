package inmem

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/models"
)

type APIKeyRepo struct {
	mu   sync.RWMutex
	data map[int]models.APIKey
	next int
}

func NewAPIKeyRepo(pool *pgxpool.Pool) *APIKeyRepo {
	return &APIKeyRepo{
		data: make(map[int]models.APIKey),
		next: 1,
	}
}

func (i *APIKeyRepo) WithTx(tx pgx.Tx) models.APIKeyRepository {
	return i
}

func (i *APIKeyRepo) List(ctx context.Context) ([]models.APIKey, error) {
	i.mu.Lock()
	defer i.mu.RUnlock()

	var out []models.APIKey
	for _, v := range i.data {
		out = append(out, v)
	}
	return out, nil
}

func (i *APIKeyRepo) Get(ctx context.Context, id int) (*models.APIKey, error) {
	i.mu.Lock()
	defer i.mu.RUnlock()

	v, ok := i.data[id]
	if !ok {
		return nil, models.ErrAPIKeyNotFound
	}

	return &v, nil
}

func (i *APIKeyRepo) Insert(ctx context.Context, a *models.APIKey) (*models.APIKey, error) {
	i.mu.Lock()
	defer i.mu.RUnlock()

	a.ID = i.next
	i.next++

	i.data[a.ID] = *a
	return a, nil
}

func (i *APIKeyRepo) Delete(ctx context.Context, id int) error {
	i.mu.Lock()
	defer i.mu.RUnlock()

	if _, ok := i.data[id]; !ok {
		return models.ErrAPIKeyNotFound
	}

	delete(i.data, id)
	return nil
}
