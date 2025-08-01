package ticketbot

import (
	"sync"
)

type InMemoryStore struct {
	tickets map[int]*Ticket
	boards  map[int]*Board
	mu      sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	t := make(map[int]*Ticket)
	b := make(map[int]*Board)
	return &InMemoryStore{
		tickets: t,
		boards:  b,
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
