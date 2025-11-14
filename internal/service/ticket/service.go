package ticket

import (
	"context"
	"sync"

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

	cwClient    *psa.Client
	ticketLocks *sync.Map
}

type FullTicket struct {
	Board   models.Board
	Ticket  models.Ticket
	Company models.Company
	Contact models.Contact
	Owner   models.Member
	Note    models.TicketNote
}

type cwData struct {
	ticket *psa.Ticket
	note   *psa.ServiceTicketNote
}

func New(b models.BoardRepository, comp models.CompanyRepository, cn models.ContactRepository,
	mem models.MemberRepository, tix models.TicketRepository, nt models.TicketNoteRepository, cl *psa.Client) *Service {
	return &Service{
		Boards:    b,
		Companies: comp,
		Contacts:  cn,
		Members:   mem,
		Tickets:   tix,
		Notes:     nt,
		cwClient:  cl,
	}
}

func (s *Service) ProcessTicket(ctx context.Context, id int) (*models.Ticket, error) {
	lock := s.getTicketLock(id)
	if !lock.TryLock() {
		lock.Lock()
	}

	defer func() {
		lock.Unlock()
	}()

}

func (s *Service) getTicketLock(id int) *sync.Mutex {
	li, _ := s.ticketLocks.LoadOrStore(id, &sync.Mutex{})
	return li.(*sync.Mutex)
}
