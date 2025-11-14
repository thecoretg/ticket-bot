package ticket

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/external/psa"
	"github.com/thecoretg/ticketbot/internal/models"
)

type Service struct {
	Boards    models.BoardRepository
	Companies models.CompanyRepository
	Contacts  models.ContactRepository
	Members   models.MemberRepository
	Tickets   models.TicketRepository
	Notes     models.TicketNoteRepository

	pool        *pgxpool.Pool
	cwClient    *psa.Client
	ticketLocks sync.Map
}

func New(pool *pgxpool.Pool, b models.BoardRepository, comp models.CompanyRepository, cn models.ContactRepository,
	mem models.MemberRepository, tix models.TicketRepository, nt models.TicketNoteRepository, cl *psa.Client) *Service {
	return &Service{
		Boards:    b,
		Companies: comp,
		Contacts:  cn,
		Members:   mem,
		Tickets:   tix,
		Notes:     nt,
		pool:      pool,
		cwClient:  cl,
	}
}

func (s *Service) withTx(tx pgx.Tx) *Service {
	return &Service{
		Boards:    s.Boards.WithTx(tx),
		Companies: s.Companies.WithTx(tx),
		Contacts:  s.Contacts.WithTx(tx),
		Members:   s.Members.WithTx(tx),
		Tickets:   s.Tickets.WithTx(tx),
		Notes:     s.Notes.WithTx(tx),
		pool:      s.pool,
		cwClient:  s.cwClient,
	}
}

func (s *Service) Run(ctx context.Context, action string, id int) (*models.FullTicket, error) {
	lock := s.getTicketLock(id)
	if !lock.TryLock() {
		lock.Lock()
	}

	defer func() {
		lock.Unlock()
	}()

	switch action {
	case "deleted":
		if err := s.Tickets.Delete(ctx, id); err != nil {
			return nil, fmt.Errorf("deleting ticket from store: %w", err)
		}
		return nil, nil
	case "added", "updated":
		t, err := s.processTicket(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("processing ticket: %w", err)
		}
		return t, nil
	default:
		return nil, fmt.Errorf("invalid action; expected 'added', 'updated', or 'deleted', got '%s'", action)
	}
}
