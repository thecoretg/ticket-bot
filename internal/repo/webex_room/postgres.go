package webexroom

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

func (p *PostgresRepo) List(ctx context.Context) ([]WebexRoom, error) {
	dbr, err := p.queries.ListWebexRooms(ctx)
	if err != nil {
		return nil, err
	}

	var r []WebexRoom
	for _, d := range dbr {
		r = append(r, roomFromPG(d))
	}

	return r, nil
}

func (p *PostgresRepo) Get(ctx context.Context, id int) (WebexRoom, error) {
	d, err := p.queries.GetWebexRoom(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return WebexRoom{}, ErrNotFound
		}
		return WebexRoom{}, err
	}

	return roomFromPG(d), nil
}

func (p *PostgresRepo) GetByWebexID(ctx context.Context, webexID string) (WebexRoom, error) {
	d, err := p.queries.GetWebexRoomByWebexID(ctx, webexID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return WebexRoom{}, ErrNotFound
		}
		return WebexRoom{}, err
	}

	return roomFromPG(d), nil
}

func (p *PostgresRepo) Upsert(ctx context.Context, r WebexRoom) (WebexRoom, error) {
	d, err := p.queries.UpsertWebexRoom(ctx, pgUpsertParams(r))
	if err != nil {
		return WebexRoom{}, err
	}

	return roomFromPG(d), nil
}

func (p *PostgresRepo) Delete(ctx context.Context, id int) error {
	if err := p.queries.DeleteWebexRoom(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func pgUpsertParams(r WebexRoom) db.UpsertWebexRoomParams {
	return db.UpsertWebexRoomParams{
		WebexID: r.WebexID,
		Name:    r.Name,
		Type:    r.Type,
	}
}

func roomFromPG(pg db.WebexRoom) WebexRoom {
	return WebexRoom{
		ID:        pg.ID,
		WebexID:   pg.WebexID,
		Name:      pg.Name,
		Type:      pg.Type,
		CreatedOn: pg.CreatedOn,
		UpdatedOn: pg.UpdatedOn,
	}
}
