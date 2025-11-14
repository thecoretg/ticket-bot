package ticket

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func TestService_getCwData(t *testing.T) {
	ctx := context.Background()
	s, err := newTestService(t, ctx)
	if err != nil {
		t.Fatalf("creating service: %v", err)
	}

	for _, id := range testTicketIDs(t) {
		cd, err := s.getCwData(id)
		if err != nil {
			t.Errorf("getting connectwise data: %v", err)
			continue
		}

		if cd.ticket.ID != id {
			t.Errorf("wanted ticket id %d, got %d", id, cd.ticket.ID)
			continue
		}
	}

}

func TestService_Run(t *testing.T) {
	ctx := context.Background()
	s, err := newTestService(t, ctx)
	if err != nil {
		t.Fatalf("creating service: %v", err)
	}

	for _, id := range testTicketIDs(t) {
		if _, err := s.Run(ctx, id); err != nil {
			t.Errorf("running service: %v", err)
			continue
		}
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

func testTicketIDs(t *testing.T) []int {
	t.Helper()
	raw := os.Getenv("TEST_TICKET_IDS")
	split := strings.Split(raw, ",")

	var ids []int
	for _, s := range split {
		i, err := strconv.Atoi(s)
		if err != nil {
			t.Logf("couldn't convert ticket id '%s' to integer", s)
			continue
		}
		ids = append(ids, i)
	}

	return ids
}
