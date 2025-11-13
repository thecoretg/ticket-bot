package userfwd

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

func (p *PostgresRepo) ListAll(ctx context.Context) ([]Forward, error) {
	dm, err := p.queries.ListWebexUserForwards(ctx)
	if err != nil {
		return nil, err
	}

	var b []Forward
	for _, d := range dm {
		b = append(b, forwardFromPG(d))
	}

	return b, nil
}

func (p *PostgresRepo) ListByEmail(ctx context.Context, email string) ([]Forward, error) {
	dm, err := p.queries.ListWebexUserForwardsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	var b []Forward
	for _, d := range dm {
		b = append(b, forwardFromPG(d))
	}

	return b, nil
}

func (p *PostgresRepo) Get(ctx context.Context, id int) (Forward, error) {
	d, err := p.queries.GetWebexUserForward(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Forward{}, ErrNotFound
		}
		return Forward{}, err
	}

	return forwardFromPG(d), nil
}

func (p *PostgresRepo) Insert(ctx context.Context, b Forward) (Forward, error) {
	d, err := p.queries.InsertWebexUserForward(ctx, pgInsertParams(b))
	if err != nil {
		return Forward{}, err
	}

	return forwardFromPG(d), nil
}

func (p *PostgresRepo) Delete(ctx context.Context, id int) error {
	if err := p.queries.DeleteWebexForward(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func pgInsertParams(t Forward) db.InsertWebexUserForwardParams {
	return db.InsertWebexUserForwardParams{
		UserEmail:     t.UserEmail,
		DestRoomID:    t.DestRoomID,
		StartDate:     t.StartDate,
		EndDate:       t.EndDate,
		Enabled:       t.Enabled,
		UserKeepsCopy: t.UserKeepsCopy,
	}
}

func forwardFromPG(pg db.WebexUserForward) Forward {
	return Forward{
		ID:            pg.ID,
		UserEmail:     pg.UserEmail,
		DestRoomID:    pg.DestRoomID,
		StartDate:     pg.StartDate,
		EndDate:       pg.EndDate,
		Enabled:       pg.Enabled,
		UserKeepsCopy: pg.UserKeepsCopy,
		UpdatedOn:     pg.UpdatedOn,
		AddedOn:       pg.CreatedOn,
	}
}
