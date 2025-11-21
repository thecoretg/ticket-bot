package inmem

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/models"
)

type WebexRoomRepo struct {
	mu   sync.RWMutex
	data map[int]models.WebexRoom
	next int
}

func NewWebexRoomRepo(pool *pgxpool.Pool) *WebexRoomRepo {
	return &WebexRoomRepo{
		data: make(map[int]models.WebexRoom),
		next: 1,
	}
}

func (p *WebexRoomRepo) WithTx(tx pgx.Tx) models.WebexRoomRepository {
	return p
}

func (p *WebexRoomRepo) List(ctx context.Context) ([]models.WebexRoom, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var out []models.WebexRoom
	for _, v := range p.data {
		out = append(out, v)
	}
	return out, nil
}

func (p *WebexRoomRepo) Get(ctx context.Context, id int) (models.WebexRoom, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	v, ok := p.data[id]
	if !ok {
		return models.WebexRoom{}, models.ErrWebexRoomNotFound
	}
	return v, nil
}

func (p *WebexRoomRepo) GetByWebexID(ctx context.Context, webexID string) (models.WebexRoom, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, v := range p.data {
		if v.WebexID == webexID {
			return v, nil
		}
	}
	return models.WebexRoom{}, models.ErrWebexRoomNotFound
}

func (p *WebexRoomRepo) Upsert(ctx context.Context, r models.WebexRoom) (models.WebexRoom, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if r.ID == 0 {
		r.ID = p.next
		p.next++
	}
	p.data[r.ID] = r
	return r, nil
}

func (p *WebexRoomRepo) Delete(ctx context.Context, id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.data[id]; !ok {
		return models.ErrWebexRoomNotFound
	}
	delete(p.data, id)
	return nil
}
