package ticketbot

import "time"

type Store interface {
	UpsertTicket(ticket *Ticket) error
	GetTicket(ticketID int) (*Ticket, error)
	ListTickets() ([]Ticket, error)

	UpsertTicketNote(ticketNote *TicketNote) error
	GetTicketNote(ticketNoteID int) (*TicketNote, error)
	ListTicketNotes(ticketID int) ([]TicketNote, error)

	UpsertBoard(board *Board) error
	GetBoard(boardID int) (*Board, error)
	ListBoards() ([]Board, error)
}

type TimeDetails struct {
	AddedToStore time.Time `json:"added_to_store"`
}

type Ticket struct {
	ID          int          `json:"id"`
	Summary     string       `json:"summary"`
	BoardID     int          `json:"board_id"`
	TicketNotes []TicketNote `json:"ticket_notes"`
	OwnerID     int          `json:"owner_id"`
	Resources   string       `json:"resources"`
	UpdatedBy   string       `json:"updated_by"`
	TimeDetails
}

type TicketNote struct {
	ID       int  `json:"id"`
	TicketID int  `json:"ticket_id"`
	Notified bool `json:"notified"`
}

type Board struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	NotifyEnabled bool   `json:"notify_enabled"`
	WebexRoomID   string `json:"webex_room_id"`
}
