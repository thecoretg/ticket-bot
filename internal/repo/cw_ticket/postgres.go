package ticket

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

func (p *PostgresRepo) List(ctx context.Context) ([]Ticket, error) {
	dm, err := p.queries.ListTickets(ctx)
	if err != nil {
		return nil, err
	}

	var b []Ticket
	for _, d := range dm {
		b = append(b, ticketFromPG(d))
	}

	return b, nil
}

func (p *PostgresRepo) Get(ctx context.Context, id int) (Ticket, error) {
	d, err := p.queries.GetTicket(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Ticket{}, ErrNotFound
		}
		return Ticket{}, err
	}

	return ticketFromPG(d), nil
}

func (p *PostgresRepo) Upsert(ctx context.Context, b Ticket) (Ticket, error) {
	d, err := p.queries.UpsertTicket(ctx, pgUpsertParams(b))
	if err != nil {
		return Ticket{}, err
	}

	return ticketFromPG(d), nil
}

func (p *PostgresRepo) Delete(ctx context.Context, id int) error {
	if err := p.queries.DeleteTicket(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func pgUpsertParams(t Ticket) db.UpsertTicketParams {
	return db.UpsertTicketParams{
		ID:        t.ID,
		Summary:   t.Summary,
		BoardID:   t.BoardID,
		OwnerID:   t.OwnerID,
		CompanyID: t.CompanyID,
		ContactID: t.ContactID,
		Resources: t.Resources,
		UpdatedBy: t.UpdatedBy,
	}
}

func ticketFromPG(pg db.CwTicket) Ticket {
	return Ticket{
		ID:        pg.ID,
		Summary:   pg.Summary,
		BoardID:   pg.BoardID,
		OwnerID:   pg.OwnerID,
		CompanyID: pg.CompanyID,
		ContactID: pg.ContactID,
		Resources: pg.Resources,
		UpdatedBy: pg.UpdatedBy,
		UpdatedOn: pg.UpdatedOn,
		AddedOn:   pg.AddedOn,
	}
}
