package inmem

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/models"
)

type CompanyRepo struct {
	mu   sync.RWMutex
	data map[int]models.Company
	next int
}

func NewCompanyRepo(pool *pgxpool.Pool) *CompanyRepo {
	return &CompanyRepo{
		data: make(map[int]models.Company),
		next: 1,
	}
}

func (p *CompanyRepo) WithTx(tx pgx.Tx) models.CompanyRepository {
	return p
}

func (p *CompanyRepo) List(ctx context.Context) ([]models.Company, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var out []models.Company
	for _, v := range p.data {
		out = append(out, v)
	}
	return out, nil
}

func (p *CompanyRepo) Get(ctx context.Context, id int) (models.Company, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	v, ok := p.data[id]
	if !ok {
		return models.Company{}, models.ErrCompanyNotFound
	}
	return v, nil
}

func (p *CompanyRepo) Upsert(ctx context.Context, c models.Company) (models.Company, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if c.ID == 0 {
		c.ID = p.next
		p.next++
	}
	p.data[c.ID] = c
	return c, nil
}

func (p *CompanyRepo) Delete(ctx context.Context, id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.data[id]; !ok {
		return models.ErrCompanyNotFound
	}
	delete(p.data, id)
	return nil
}
