package inmem

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/models"
)

type ConfigRepo struct {
	mu   sync.RWMutex
	data map[int]models.Config
	next int
}

func NewConfigRepo(pool *pgxpool.Pool) *ConfigRepo {
	return &ConfigRepo{
		data: make(map[int]models.Config),
		next: 1,
	}
}

func (p *ConfigRepo) WithTx(tx pgx.Tx) models.ConfigRepository {
	return p
}

func (p *ConfigRepo) Get(ctx context.Context) (*models.Config, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// return the first config if exists
	for _, v := range p.data {
		return &v, nil
	}
	return nil, models.ErrConfigNotFound
}

func (p *ConfigRepo) InsertDefault(ctx context.Context) (*models.Config, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	c := models.Config{
		ID: p.next,
	}
	p.next++
	p.data[c.ID] = c
	return &c, nil
}

func (p *ConfigRepo) Upsert(ctx context.Context, c *models.Config) (*models.Config, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if c.ID == 0 {
		c.ID = p.next
		p.next++
	}
	p.data[c.ID] = *c
	return c, nil
}
