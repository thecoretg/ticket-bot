package ticket

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/thecoretg/ticketbot/internal/external/psa"
	"github.com/thecoretg/ticketbot/internal/postgres"
)

const testTicketID = 698014

func TestNewService(t *testing.T) {
	if _, err := newTestService(t, context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestService_getCwData(t *testing.T) {
	ctx := context.Background()
	s, err := newTestService(t, ctx)
	if err != nil {
		t.Fatalf("creating service: %v", err)
	}

	cd, err := s.getCwData(testTicketID)
	if err != nil {
		t.Fatal(err)
	}

	if cd.ticket.ID != testTicketID {
		t.Fatalf("wanted ticket id %d, got %d", testTicketID, cd.ticket.ID)
	}
}

func TestService_Run(t *testing.T) {
	ctx := context.Background()
	s, err := newTestService(t, ctx)
	if err != nil {
		t.Fatalf("creating service: %v", err)
	}

	if _, err := s.Run(ctx, testTicketID); err != nil {
		t.Fatalf("running service: %v", err)
	}
}

func newTestService(t *testing.T, ctx context.Context) (*Service, error) {
	t.Helper()
	if err := godotenv.Load("testing.env"); err != nil {
		return nil, fmt.Errorf("loading .env")
	}

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		return nil, errors.New("postgres dsn is empty")
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("creating pgx pool: %w", err)
	}

	t.Cleanup(func() { pool.Close() })

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pinging pool")
	}

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
