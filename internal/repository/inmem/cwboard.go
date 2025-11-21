package inmem

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/models"
)

type BoardRepo struct {
	mu   sync.RWMutex
	data map[int]models.Board
	next int
}

func NewBoardRepo(pool *pgxpool.Pool) *BoardRepo {
	return &BoardRepo{
		data: make(map[int]models.Board),
		next: 1,
	}
}

func (p *BoardRepo) WithTx(tx pgx.Tx) models.BoardRepository {
	return p
}

func (p *BoardRepo) List(ctx context.Context) ([]models.Board, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var out []models.Board
	for _, v := range p.data {
		out = append(out, v)
	}
	return out, nil
}

func (p *BoardRepo) Get(ctx context.Context, id int) (models.Board, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	v, ok := p.data[id]
	if !ok {
		return models.Board{}, models.ErrBoardNotFound
	}
	return v, nil
}

func (p *BoardRepo) Upsert(ctx context.Context, b models.Board) (models.Board, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if b.ID == 0 {
		b.ID = p.next
		p.next++
	}
	p.data[b.ID] = b
	return b, nil
}

func (p *BoardRepo) Delete(ctx context.Context, id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.data[id]; !ok {
		return models.ErrBoardNotFound
	}
	delete(p.data, id)
	return nil
}
