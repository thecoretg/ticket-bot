package ticketbot

import (
	"sync"
)

type InMemoryStore struct {
	tickets map[int]*Ticket
	boards  map[int]*Board
	users   map[int]*User
	apiKeys map[int]*APIKey
	mu      sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	t := make(map[int]*Ticket)
	b := make(map[int]*Board)
	u := make(map[int]*User)
	a := make(map[int]*APIKey)
	return &InMemoryStore{
		tickets: t,
		boards:  b,
		users:   u,
		apiKeys: a,
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

func (m *InMemoryStore) UpsertUser(user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.ID] = user
	return nil
}

func (m *InMemoryStore) GetUser(userID int) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if user, exists := m.users[userID]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *InMemoryStore) ListUsers() ([]User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var users []User
	for _, user := range m.users {
		users = append(users, *user)
	}
	if users == nil {
		users = []User{}
	}
	return users, nil
}

func (m *InMemoryStore) UpsertAPIKey(apiKey *APIKey) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.apiKeys[apiKey.ID] = apiKey
	return nil
}

func (m *InMemoryStore) GetAPIKey(apiKeyID int) (*APIKey, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if apiKey, exists := m.apiKeys[apiKeyID]; exists {
		return apiKey, nil
	}
	return nil, nil
}

func (m *InMemoryStore) ListAPIKeys() ([]APIKey, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var apiKeys []APIKey
	for _, apiKey := range m.apiKeys {
		apiKeys = append(apiKeys, *apiKey)
	}
	if apiKeys == nil {
		apiKeys = []APIKey{}
	}
	return apiKeys, nil
}
