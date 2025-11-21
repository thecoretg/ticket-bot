package inmem

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/models"
)

type NotifierRepo struct {
	mu   sync.RWMutex
	data map[int]models.Notifier
	next int
}

func NewNotifierRepo(pool *pgxpool.Pool) *NotifierRepo {
	return &NotifierRepo{
		data: make(map[int]models.Notifier),
		next: 1,
	}
}

func (p *NotifierRepo) WithTx(tx pgx.Tx) models.NotifierRepository {
	return p
}

func (p *NotifierRepo) ListAll(ctx context.Context) ([]models.Notifier, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var out []models.Notifier
	for _, v := range p.data {
		out = append(out, v)
	}
	return out, nil
}

func (p *NotifierRepo) ListByBoard(ctx context.Context, boardID int) ([]models.Notifier, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var out []models.Notifier
	for _, v := range p.data {
		if v.CwBoardID == boardID {
			out = append(out, v)
		}
	}
	return out, nil
}

func (p *NotifierRepo) ListByRoom(ctx context.Context, roomID int) ([]models.Notifier, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var out []models.Notifier
	for _, v := range p.data {
		if v.WebexRoomID == roomID {
			out = append(out, v)
		}
	}
	return out, nil
}

func (p *NotifierRepo) Get(ctx context.Context, id int) (*models.Notifier, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	v, ok := p.data[id]
	if !ok {
		return nil, models.ErrNotifierNotFound
	}
	return &v, nil
}

func (p *NotifierRepo) Exists(ctx context.Context, boardID, roomID int) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, v := range p.data {
		if v.CwBoardID == boardID && v.WebexRoomID == roomID {
			return true, nil
		}
	}
	return false, nil
}

func (p *NotifierRepo) Insert(ctx context.Context, n *models.Notifier) (*models.Notifier, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	n.ID = p.next
	p.next++
	p.data[n.ID] = *n
	return n, nil
}

func (p *NotifierRepo) Update(ctx context.Context, n *models.Notifier) (*models.Notifier, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.data[n.ID]; !ok {
		return nil, models.ErrNotifierNotFound
	}
	p.data[n.ID] = *n
	return n, nil
}

func (p *NotifierRepo) Delete(ctx context.Context, id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.data[id]; !ok {
		return models.ErrNotifierNotFound
	}
	delete(p.data, id)
	return nil
}
