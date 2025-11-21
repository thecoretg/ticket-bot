package inmem

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/models"
)

type TicketRepo struct {
	mu   sync.RWMutex
	data map[int]models.Ticket
	next int
}

func NewTicketRepo(pool *pgxpool.Pool) *TicketRepo {
	return &TicketRepo{
		data: make(map[int]models.Ticket),
		next: 1,
	}
}

func (p *TicketRepo) WithTx(tx pgx.Tx) models.TicketRepository {
	return p
}

func (p *TicketRepo) List(ctx context.Context) ([]models.Ticket, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var out []models.Ticket
	for _, v := range p.data {
		out = append(out, v)
	}
	return out, nil
}

func (p *TicketRepo) Get(ctx context.Context, id int) (models.Ticket, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	v, ok := p.data[id]
	if !ok {
		return models.Ticket{}, models.ErrTicketNotFound
	}
	return v, nil
}

func (p *TicketRepo) Upsert(ctx context.Context, t models.Ticket) (models.Ticket, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if t.ID == 0 {
		t.ID = p.next
		p.next++
	}
	p.data[t.ID] = t
	return t, nil
}

func (p *TicketRepo) Delete(ctx context.Context, id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.data[id]; !ok {
		return models.ErrTicketNotFound
	}
	delete(p.data, id)
	return nil
}
