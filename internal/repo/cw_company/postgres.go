package company

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

func (p *PostgresRepo) List(ctx context.Context) ([]Company, error) {
	dbs, err := p.queries.ListCompanies(ctx)
	if err != nil {
		return nil, err
	}

	var b []Company
	for _, d := range dbs {
		b = append(b, companyFromPG(d))
	}

	return b, nil
}

func (p *PostgresRepo) Get(ctx context.Context, id int) (Company, error) {
	d, err := p.queries.GetCompany(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Company{}, ErrNotFound
		}
		return Company{}, err
	}

	return companyFromPG(d), nil
}

func (p *PostgresRepo) Upsert(ctx context.Context, b Company) (Company, error) {
	d, err := p.queries.UpsertCompany(ctx, pgUpsertParams(b))
	if err != nil {
		return Company{}, err
	}

	return companyFromPG(d), nil
}

func (p *PostgresRepo) Delete(ctx context.Context, id int) error {
	if err := p.queries.DeleteCompany(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func pgUpsertParams(c Company) db.UpsertCompanyParams {
	return db.UpsertCompanyParams{
		ID:   c.ID,
		Name: c.Name,
	}
}

func companyFromPG(pg db.CwCompany) Company {
	return Company{
		ID:        pg.ID,
		Name:      pg.Name,
		UpdatedOn: pg.UpdatedOn,
		AddedOn:   pg.AddedOn,
	}
}
