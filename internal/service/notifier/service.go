package notifier

import (
	"github.com/thecoretg/ticketbot/internal/models"
)

type Service struct {
	Config    models.ConfigRepository
	Boards    models.BoardRepository
	Companies models.CompanyRepository
	Contacts  models.ContactRepository
	Members   models.MemberRepository
	Tickets   models.TicketRepository
	Notes     models.TicketNoteRepository
	Rooms     models.WebexRoomRepository
	Notifiers models.NotifierRepository
	Forwards  models.UserForwardRepository
}
