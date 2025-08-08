package ticketbot

import (
	"sync"
)

type InMemoryStore struct {
	tickets     map[int]*Ticket
	boards      map[int]*Board
	ticketNotes map[int]*TicketNote
	mu          sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	t := make(map[int]*Ticket)
	b := make(map[int]*Board)
	n := make(map[int]*TicketNote)
	return &InMemoryStore{
		tickets:     t,
		boards:      b,
		ticketNotes: n,
	}
}

func (m *InMemoryStore) UpsertTicket(ticket *Ticket) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tickets[ticket.ID] = ticket
	return nil
}

func (m *InMemoryStore) GetTicket(ticketID int) (*Ticket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if ticket, exists := m.tickets[ticketID]; exists {
		return ticket, nil
	}
	return nil, nil
}

func (m *InMemoryStore) ListTickets() ([]Ticket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var tickets []Ticket
	for _, ticket := range m.tickets {
		tickets = append(tickets, *ticket)
	}

	if tickets == nil {
		tickets = []Ticket{}
	}

	return tickets, nil
}

func (m *InMemoryStore) UpsertBoard(board *Board) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.boards[board.ID] = board
	return nil
}

func (m *InMemoryStore) GetBoard(boardID int) (*Board, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if board, exists := m.boards[boardID]; exists {
		return board, nil
	}
	return nil, nil
}

func (m *InMemoryStore) ListBoards() ([]Board, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var boards []Board
	for _, board := range m.boards {
		boards = append(boards, *board)
	}

	if boards == nil {
		boards = []Board{}
	}

	return boards, nil
}

func (m *InMemoryStore) UpsertTicketNote(ticketNote *TicketNote) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ticketNotes[ticketNote.ID] = ticketNote
	return nil
}

func (m *InMemoryStore) GetTicketNote(ticketNoteID int) (*TicketNote, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if note, exists := m.ticketNotes[ticketNoteID]; exists {
		return note, nil
	}
	return nil, nil
}

func (m *InMemoryStore) ListTicketNotes(ticketID int) ([]TicketNote, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var notes []TicketNote
	for _, note := range m.ticketNotes {
		if note.TicketID == ticketID {
			notes = append(notes, *note)
		}
	}
	if notes == nil {
		notes = []TicketNote{}
	}
	return notes, nil
}
