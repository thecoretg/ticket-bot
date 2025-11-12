package board

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

func NewPostgresRepo(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (p *PostgresRepo) List(ctx context.Context) ([]Board, error) {
	dbs, err := p.queries.ListBoards(ctx)
	if err != nil {
		return nil, err
	}

	var b []Board
	for _, d := range dbs {
		b = append(b, boardFromPG(d))
	}

	return b, nil
}

func (p *PostgresRepo) Get(ctx context.Context, id int) (Board, error) {
	d, err := p.queries.GetBoard(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Board{}, ErrNotFound
		}
		return Board{}, err
	}

	return boardFromPG(d), nil
}

func (p *PostgresRepo) Upsert(ctx context.Context, b Board) (Board, error) {
	d, err := p.queries.UpsertBoard(ctx, pgUpsertParams(b))
	if err != nil {
		return Board{}, err
	}

	return boardFromPG(d), nil
}

func (p *PostgresRepo) Delete(ctx context.Context, id int) error {
	if err := p.queries.DeleteBoard(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func pgUpsertParams(b Board) db.UpsertBoardParams {
	return db.UpsertBoardParams{
		ID:   b.ID,
		Name: b.Name,
	}
}

func boardFromPG(pg db.CwBoard) Board {
	return Board{
		ID:        pg.ID,
		Name:      pg.Name,
		UpdatedOn: pg.UpdatedOn,
		AddedOn:   pg.AddedOn,
	}
}
