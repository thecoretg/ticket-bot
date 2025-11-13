package cwcontact

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

func (p *PostgresRepo) List(ctx context.Context) ([]Contact, error) {
	dbs, err := p.queries.ListContacts(ctx)
	if err != nil {
		return nil, err
	}

	var b []Contact
	for _, d := range dbs {
		b = append(b, contactFromPG(d))
	}

	return b, nil
}

func (p *PostgresRepo) Get(ctx context.Context, id int) (Contact, error) {
	d, err := p.queries.GetContact(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Contact{}, ErrNotFound
		}
		return Contact{}, err
	}

	return contactFromPG(d), nil
}

func (p *PostgresRepo) Upsert(ctx context.Context, b Contact) (Contact, error) {
	d, err := p.queries.UpsertContact(ctx, pgUpsertParams(b))
	if err != nil {
		return Contact{}, err
	}

	return contactFromPG(d), nil
}

func (p *PostgresRepo) Delete(ctx context.Context, id int) error {
	if err := p.queries.DeleteContact(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func pgUpsertParams(c Contact) db.UpsertContactParams {
	return db.UpsertContactParams{
		ID:        c.ID,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		CompanyID: c.CompanyID,
	}
}

func contactFromPG(pg db.CwContact) Contact {
	return Contact{
		ID:        pg.ID,
		FirstName: pg.FirstName,
		LastName:  pg.LastName,
		CompanyID: pg.CompanyID,
		UpdatedOn: pg.UpdatedOn,
		AddedOn:   pg.AddedOn,
	}
}
