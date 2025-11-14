package ticket

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/thecoretg/ticketbot/internal/external/psa"
	"github.com/thecoretg/ticketbot/internal/postgres"
)

func TestNewService(t *testing.T) {
	if _, err := newTestService(t, context.Background()); err != nil {
		t.Fatal(err)
	}
}

func newTestService(t *testing.T, ctx context.Context) (*Service, error) {
	t.Helper()
	if err := godotenv.Load("testing.env"); err != nil {
		return nil, fmt.Errorf("loading .env")
	}

	dsn := os.Getenv("POSTGRES_DSN")
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("creating pgx pool: %w", err)
	}

	t.Cleanup(func() { pool.Close() })

	cwCreds := &psa.Creds{
		PublicKey:  os.Getenv("CW_PUB_KEY"),
		PrivateKey: os.Getenv("CW_PRIV_KEY"),
		ClientId:   os.Getenv("CW_CLIENT_ID"),
		CompanyId:  os.Getenv("CW_COMPANY_ID"),
	}

	return New(
		pool,
		postgres.NewBoardRepo(pool),
		postgres.NewCompanyRepo(pool),
		postgres.NewContactRepo(pool),
		postgres.NewMemberRepo(pool),
		postgres.NewTicketRepo(pool),
		postgres.NewTicketNoteRepo(pool),
		psa.NewClient(cwCreds),
	), nil
}
