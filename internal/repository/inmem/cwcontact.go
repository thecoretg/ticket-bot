package inmem

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/models"
)

type ContactRepo struct {
	mu   sync.RWMutex
	data map[int]models.Contact
	next int
}

func NewContactRepo(pool *pgxpool.Pool) *ContactRepo {
	return &ContactRepo{
		data: make(map[int]models.Contact),
		next: 1,
	}
}

func (p *ContactRepo) WithTx(tx pgx.Tx) models.ContactRepository {
	return p
}

func (p *ContactRepo) List(ctx context.Context) ([]models.Contact, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var out []models.Contact
	for _, v := range p.data {
		out = append(out, v)
	}
	return out, nil
}

func (p *ContactRepo) Get(ctx context.Context, id int) (models.Contact, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	v, ok := p.data[id]
	if !ok {
		return models.Contact{}, models.ErrContactNotFound
	}
	return v, nil
}

func (p *ContactRepo) Upsert(ctx context.Context, c models.Contact) (models.Contact, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if c.ID == 0 {
		c.ID = p.next
		p.next++
	}
	p.data[c.ID] = c
	return c, nil
}

func (p *ContactRepo) Delete(ctx context.Context, id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.data[id]; !ok {
		return models.ErrContactNotFound
	}
	delete(p.data, id)
	return nil
}
