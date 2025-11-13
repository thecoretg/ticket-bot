package notifier

import (
	"github.com/thecoretg/ticketbot/internal/repo/config"
	board "github.com/thecoretg/ticketbot/internal/repo/cw_board"
	company "github.com/thecoretg/ticketbot/internal/repo/cw_company"
	contact "github.com/thecoretg/ticketbot/internal/repo/cw_contact"
	member "github.com/thecoretg/ticketbot/internal/repo/cw_member"
	ticket "github.com/thecoretg/ticketbot/internal/repo/cw_ticket"
	note "github.com/thecoretg/ticketbot/internal/repo/cw_ticket_note"
	notifier "github.com/thecoretg/ticketbot/internal/repo/msg_notifier"
	userfwd "github.com/thecoretg/ticketbot/internal/repo/msg_userfwd"
	webexroom "github.com/thecoretg/ticketbot/internal/repo/webex_room"
)

type Service struct {
	Config    config.Repository
	Boards    board.Repository
	Companies company.Repository
	Contacts  contact.Repository
	Members   member.Repository
	Tickets   ticket.Repository
	Notes     note.Repository
	Rooms     webexroom.Repository
	Notifiers notifier.Repository
	Forwards  userfwd.Repository
}
