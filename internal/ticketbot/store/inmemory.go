package store

import (
	"errors"
	"fmt"
	"tctg-automation/internal/ticketbot/types"
)

type InMemoryStore struct {
	store map[int]*types.Ticket
}

func NewInMemoryStore() *InMemoryStore {
	s := make(map[int]*types.Ticket)
	return &InMemoryStore{
		store: s,
	}
}

func (m *InMemoryStore) AddTicket(ticket *types.Ticket) error {
	if ticket == nil {
		return errors.New("ticket cannot be nil")
	}

	if _, exists := m.store[ticket.ID]; exists {
		return fmt.Errorf("ticket with ID %d already exists", ticket.ID)
	}

	m.store[ticket.ID] = ticket
	return nil
}

func (m *InMemoryStore) UpdateTicket(ticket *types.Ticket) error {
	if ticket == nil {
		return errors.New("ticket cannot be nil")
	}

	if _, exists := m.store[ticket.ID]; !exists {
		return fmt.Errorf("ticket with ID %d does not exist", ticket.ID)
	}

	m.store[ticket.ID] = ticket
	return nil
}

func (m *InMemoryStore) GetTicket(ticketID int) (*types.Ticket, error) {
	if ticket, exists := m.store[ticketID]; exists {
		return ticket, nil
	}
	return nil, fmt.Errorf("ticket with ID %d not found", ticketID)
}

func (m *InMemoryStore) ListTickets() ([]types.Ticket, error) {
	var tickets []types.Ticket
	for _, ticket := range m.store {
		tickets = append(tickets, *ticket)
	}
	return tickets, nil
}
