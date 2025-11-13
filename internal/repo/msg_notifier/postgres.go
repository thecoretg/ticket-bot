package notifier

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/db"
)

type PostgresRepo struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresRepo(pool *pgxpool.Pool, q *db.Queries) *PostgresRepo {
	return &PostgresRepo{
		pool:    pool,
		queries: q,
	}
}

func (p *PostgresRepo) ListAll(ctx context.Context) ([]Notifier, error) {
	dm, err := p.queries.ListNotifierConnections(ctx)
	if err != nil {
		return nil, err
	}

	var b []Notifier
	for _, d := range dm {
		b = append(b, notifierFromPG(d))
	}

	return b, nil
}

func (p *PostgresRepo) ListByBoard(ctx context.Context, boardID int) ([]Notifier, error) {
	dm, err := p.queries.ListNotifierConnectionsByBoard(ctx, boardID)
	if err != nil {
		return nil, err
	}

	var b []Notifier
	for _, d := range dm {
		b = append(b, notifierFromPG(d))
	}

	return b, nil
}

func (p *PostgresRepo) ListByRoom(ctx context.Context, roomID int) ([]Notifier, error) {
	dm, err := p.queries.ListNotifierConnectionsByRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}

	var b []Notifier
	for _, d := range dm {
		b = append(b, notifierFromPG(d))
	}

	return b, nil
}

func (p *PostgresRepo) Get(ctx context.Context, id int) (Notifier, error) {
	d, err := p.queries.GetNotifierConnection(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Notifier{}, ErrNotFound
		}
		return Notifier{}, err
	}

	return notifierFromPG(d), nil
}

func (p *PostgresRepo) Insert(ctx context.Context, n Notifier) (Notifier, error) {
	d, err := p.queries.InsertNotifierConnection(ctx, pgInsertParams(n))
	if err != nil {
		return Notifier{}, err
	}

	return notifierFromPG(d), nil
}

func (p *PostgresRepo) Update(ctx context.Context, n Notifier) (Notifier, error) {
	d, err := p.queries.UpdateNotifierConnection(ctx, pgUpdateParams(n))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Notifier{}, ErrNotFound
		}
		return Notifier{}, err
	}

	return notifierFromPG(d), nil
}

func (p *PostgresRepo) Delete(ctx context.Context, id int) error {
	if err := p.queries.DeleteNotifierConnection(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func pgInsertParams(n Notifier) db.InsertNotifierConnectionParams {
	return db.InsertNotifierConnectionParams{
		CwBoardID:     n.CwBoardID,
		WebexRoomID:   n.WebexRoomID,
		NotifyEnabled: n.NotifyEnabled,
	}
}

func pgUpdateParams(n Notifier) db.UpdateNotifierConnectionParams {
	return db.UpdateNotifierConnectionParams{
		ID:            n.ID,
		CwBoardID:     n.CwBoardID,
		WebexRoomID:   n.WebexRoomID,
		NotifyEnabled: n.NotifyEnabled,
	}
}

func notifierFromPG(pg db.NotifierConnection) Notifier {
	return Notifier{
		ID:            pg.ID,
		CwBoardID:     pg.CwBoardID,
		WebexRoomID:   pg.WebexRoomID,
		NotifyEnabled: pg.NotifyEnabled,
		CreatedOn:     pg.CreatedOn,
	}
}
