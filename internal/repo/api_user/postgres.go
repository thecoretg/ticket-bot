package apiuser

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

func (p *PostgresRepo) List(ctx context.Context) ([]APIUser, error) {
	dk, err := p.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	var k []APIUser
	for _, d := range dk {
		k = append(k, userFromPG(d))
	}

	return k, nil
}

func (p *PostgresRepo) Get(ctx context.Context, id int) (APIUser, error) {
	d, err := p.queries.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return APIUser{}, ErrNotFound
		}
		return APIUser{}, ErrNotFound
	}

	return userFromPG(d), nil
}

func (p *PostgresRepo) GetByEmail(ctx context.Context, email string) (APIUser, error) {
	d, err := p.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return APIUser{}, ErrNotFound
		}
		return APIUser{}, err
	}

	return userFromPG(d), nil
}

func (p *PostgresRepo) Insert(ctx context.Context, email string) (APIUser, error) {
	d, err := p.queries.InsertUser(ctx, email)
	if err != nil {
		return APIUser{}, err
	}

	return userFromPG(d), nil
}

func (p *PostgresRepo) Update(ctx context.Context, u APIUser) (APIUser, error) {
	d, err := p.queries.UpdateUser(ctx, pgUpdateParams(u))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return APIUser{}, ErrNotFound
		}
		return APIUser{}, err
	}

	return userFromPG(d), nil
}

func (p *PostgresRepo) Delete(ctx context.Context, id int) error {
	if err := p.queries.DeleteUser(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func pgUpdateParams(u APIUser) db.UpdateUserParams {
	return db.UpdateUserParams{
		ID:           u.ID,
		EmailAddress: u.EmailAddress,
	}
}

func userFromPG(pg db.ApiUser) APIUser {
	return APIUser{
		ID:           pg.ID,
		EmailAddress: pg.EmailAddress,
		CreatedOn:    pg.CreatedOn,
		UpdatedOn:    pg.UpdatedOn,
	}
}
