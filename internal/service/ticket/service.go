package ticket

import (
	"context"
	"sync"

	"github.com/thecoretg/ticketbot/internal/external/psa"
	cwboard "github.com/thecoretg/ticketbot/internal/repo/cw_board"
	cwcompany "github.com/thecoretg/ticketbot/internal/repo/cw_company"
	cwcontact "github.com/thecoretg/ticketbot/internal/repo/cw_contact"
	cwmember "github.com/thecoretg/ticketbot/internal/repo/cw_member"
	cwticket "github.com/thecoretg/ticketbot/internal/repo/cw_ticket"
	cwnote "github.com/thecoretg/ticketbot/internal/repo/cw_ticket_note"
)

type Service struct {
	Boards    cwboard.Repository
	Companies cwcompany.Repository
	Contacts  cwcontact.Repository
	Members   cwmember.Repository
	Tickets   cwticket.Repository
	Notes     cwnote.Repository

	cwClient    *psa.Client
	ticketLocks *sync.Map
}

type FullTicket struct {
	Board   cwboard.Board
	Ticket  cwticket.Ticket
	Company cwcompany.Company
	Contact cwcontact.Contact
	Owner   cwmember.Member
	Note    cwnote.TicketNote
}

type cwData struct {
	ticket *psa.Ticket
	note   *psa.ServiceTicketNote
}

func New(b cwboard.Repository, comp cwcompany.Repository, cn cwcontact.Repository,
	mem cwmember.Repository, tix cwticket.Repository, nt cwnote.Repository, cl *psa.Client) *Service {
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

func (s *Service) ProcessTicket(ctx context.Context, id int) (*cwticket.Ticket, error) {
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
