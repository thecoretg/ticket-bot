package inmem

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/models"
)

type TicketNoteRepo struct {
	mu   sync.RWMutex
	data map[int]models.TicketNote
	next int
}

func NewTicketNoteRepo(pool *pgxpool.Pool) *TicketNoteRepo {
	return &TicketNoteRepo{
		data: make(map[int]models.TicketNote),
		next: 1,
	}
}

func (p *TicketNoteRepo) WithTx(tx pgx.Tx) models.TicketNoteRepository {
	return p
}

func (p *TicketNoteRepo) ListByTicketID(ctx context.Context, ticketID int) ([]models.TicketNote, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var out []models.TicketNote
	for _, v := range p.data {
		if v.TicketID == ticketID {
			out = append(out, v)
		}
	}
	return out, nil
}

func (p *TicketNoteRepo) ListAll(ctx context.Context) ([]models.TicketNote, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var out []models.TicketNote
	for _, v := range p.data {
		out = append(out, v)
	}
	return out, nil
}

func (p *TicketNoteRepo) Get(ctx context.Context, id int) (models.TicketNote, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	v, ok := p.data[id]
	if !ok {
		return models.TicketNote{}, models.ErrTicketNoteNotFound
	}
	return v, nil
}

func (p *TicketNoteRepo) Upsert(ctx context.Context, n models.TicketNote) (models.TicketNote, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if n.ID == 0 {
		n.ID = p.next
		p.next++
	}
	p.data[n.ID] = n
	return n, nil
}

func (p *TicketNoteRepo) Delete(ctx context.Context, id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.data[id]; !ok {
		return models.ErrTicketNoteNotFound
	}
	delete(p.data, id)
	return nil
}
