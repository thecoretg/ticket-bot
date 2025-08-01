package ticketbot

import (
	"fmt"
	"log/slog"
	"sync"
	"tctg-automation/internal/ticketbot/types"
	"tctg-automation/pkg/connectwise"
)

const maxConcurrentPreload = 10

func (s *server) preloadFromConnectwise(preloadBoards, preloadTickets bool) error {
	if preloadBoards {
		if err := s.preloadBoards(); err != nil {
			return fmt.Errorf("preloading active boards: %w", err)
		}
	}

	if preloadTickets {
		if err := s.preloadOpenTickets(); err != nil {
			return fmt.Errorf("preloading open tickets: %w", err)
		}
	}

	return nil
}

func (s *server) preloadBoards() error {
	params := map[string]string{
		"conditions": "inactiveFlag = false",
	}

	slog.Info("loading existing boards")
	boards, err := s.cwClient.ListBoards(params)
	if err != nil {
		return fmt.Errorf("getting boards from CW: %w", err)
	}
	slog.Info("got boards", "total_boards", len(boards))
	sem := make(chan struct{}, maxConcurrentPreload)
	var wg sync.WaitGroup
	for _, board := range boards {
		sem <- struct{}{}
		wg.Add(1)
		go func(board connectwise.Board) {
			defer wg.Done()
			defer func() { <-sem }()
			b := &types.Board{
				ID:            board.ID,
				Name:          board.Name,
				NotifyEnabled: false,
				WebexRoomIDs:  nil,
			}
			if err := s.dataStore.UpsertBoard(b); err != nil {
				slog.Warn("error preloading board", "board_id", board.ID, "error", err)
			}
			slog.Info("preloaded board", "board_id", board.ID, "board_name", board.Name)
		}(board)
	}

	wg.Wait()
	return nil
}

func (s *server) preloadOpenTickets() error {
	params := map[string]string{
		"pageSize":   "100",
		"conditions": "closedFlag = false and board/id = 34",
	}

	slog.Info("loading existing open tickets")
	openTickets, err := s.cwClient.ListTickets(params)
	if err != nil {
		return fmt.Errorf("getting open tickets from CW: %w", err)
	}
	slog.Info("got open tickets", "total_tickets", len(openTickets))
	sem := make(chan struct{}, maxConcurrentPreload)
	var wg sync.WaitGroup

	for _, ticket := range openTickets {
		sem <- struct{}{}
		wg.Add(1)
		go func(ticket connectwise.Ticket) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := s.addOrUpdateTicket(nil, &ticket, false); err != nil {
				slog.Warn("error preloading open ticket", "ticket_id", ticket.ID, "error", err)
			}
		}(ticket)
	}

	wg.Wait()
	return nil
}
