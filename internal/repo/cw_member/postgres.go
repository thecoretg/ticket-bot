package member

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

func (p *PostgresRepo) List(ctx context.Context) ([]Member, error) {
	dm, err := p.queries.ListMembers(ctx)
	if err != nil {
		return nil, err
	}

	var b []Member
	for _, d := range dm {
		b = append(b, memberFromPG(d))
	}

	return b, nil
}

func (p *PostgresRepo) Get(ctx context.Context, id int) (Member, error) {
	d, err := p.queries.GetMember(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Member{}, ErrNotFound
		}
		return Member{}, err
	}

	return memberFromPG(d), nil
}

func (p *PostgresRepo) Upsert(ctx context.Context, b Member) (Member, error) {
	d, err := p.queries.UpsertMember(ctx, pgUpsertParams(b))
	if err != nil {
		return Member{}, err
	}

	return memberFromPG(d), nil
}

func (p *PostgresRepo) Delete(ctx context.Context, id int) error {
	if err := p.queries.DeleteMember(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func pgUpsertParams(m Member) db.UpsertMemberParams {
	return db.UpsertMemberParams{
		ID:           m.ID,
		Identifier:   m.Identifier,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		PrimaryEmail: m.PrimaryEmail,
	}
}

func memberFromPG(pg db.CwMember) Member {
	return Member{
		ID:           pg.ID,
		Identifier:   pg.Identifier,
		FirstName:    pg.FirstName,
		LastName:     pg.LastName,
		PrimaryEmail: pg.PrimaryEmail,
		UpdatedOn:    pg.UpdatedOn,
		AddedOn:      pg.AddedOn,
	}
}
