package config

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

func (p *PostgresRepo) Get(ctx context.Context) (Config, error) {
	d, err := p.queries.GetAppConfig(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Config{}, ErrNotFound
		}
		return Config{}, err
	}

	return configFromPG(d), nil
}

func (p *PostgresRepo) InsertDefault(ctx context.Context) (Config, error) {
	d, err := p.queries.InsertDefaultAppConfig(ctx)
	if err != nil {
		return Config{}, err
	}

	return configFromPG(d), nil
}

func (p *PostgresRepo) Upsert(ctx context.Context, c Config) (Config, error) {
	d, err := p.queries.UpsertAppConfig(ctx, pgUpsertParams(c))
	if err != nil {
		return Config{}, err
	}

	return configFromPG(d), nil
}

func NewPostgresRepo(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{
		pool:    pool,
		queries: db.New(pool),
	}
}

func pgUpsertParams(c Config) db.UpsertAppConfigParams {
	return db.UpsertAppConfigParams{
		Debug:              c.Debug,
		AttemptNotify:      c.AttemptNotify,
		MaxMessageLength:   c.MaxMessageLength,
		MaxConcurrentSyncs: c.MaxConcurrentSyncs,
	}
}

func configFromPG(pg db.AppConfig) Config {
	return Config{
		ID:                 pg.ID,
		Debug:              pg.Debug,
		AttemptNotify:      pg.AttemptNotify,
		MaxMessageLength:   pg.MaxMessageLength,
		MaxConcurrentSyncs: pg.MaxConcurrentSyncs,
	}
}
