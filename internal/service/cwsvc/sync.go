package cwsvc

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/thecoretg/ticketbot/internal/external/psa"
	"github.com/thecoretg/ticketbot/internal/models"
)

func (s *Service) SyncOpenTickets(ctx context.Context, boardIDs []int, maxSyncs int) error {
	start := time.Now()
	slog.Debug("beginning ticket sync", "board_ids", boardIDs)
	con := "closedFlag = false"
	if len(boardIDs) > 0 {
		con += fmt.Sprintf(" AND %s", boardIDParam(boardIDs))
	}

	params := map[string]string{
		"pageSize":   "100",
		"conditions": con,
	}

	tix, err := s.cwClient.ListTickets(params)
	if err != nil {
		return fmt.Errorf("getting open tickets from connectwise: %w", err)
	}
	slog.Debug("open ticket sync: got open tickets from connectwise", "total_tickets", len(tix))
	sem := make(chan struct{}, maxSyncs)
	var wg sync.WaitGroup
	errCh := make(chan error, len(tix))

	for _, t := range tix {
		sem <- struct{}{}
		wg.Add(1)
		go func(ticket psa.Ticket) {
			defer func() { <-sem }()
			if _, err := s.processTicket(ctx, t.ID); err != nil {
				slog.Error("ticket sync error", "ticket_id", t.ID, "error", err)
				errCh <- fmt.Errorf("error syncing ticket %d: %w", t.ID, err)
			} else {
				errCh <- nil
			}
		}(t)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			slog.Error("syncing ticket", "error", err)
		}
	}
	slog.Info("syncing tickets complete", "took_time", time.Since(start))
	return nil
}

func (s *Service) SyncBoards(ctx context.Context) error {
	start := time.Now()
	slog.Debug("beginning connectwise board sync")
	cwb, err := s.cwClient.ListBoards(nil)
	if err != nil {
		return fmt.Errorf("listing connectwise boards: %w", err)
	}
	slog.Debug("board sync: got boards from connectwise", "total_boards", len(cwb))
	boards := make(map[int]psa.Board, len(cwb))
	for _, b := range boards {
		boards[b.ID] = b
	}

	sb, err := s.Boards.List(ctx)
	if err != nil {
		return fmt.Errorf("listing boards from store: %w", err)
	}

	storeByID := make(map[int]models.Board, len(sb))
	for _, b := range sb {
		storeByID[b.ID] = b
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("beginning tx: %w", err)
	}

	txSvc := s.withTx(tx)

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// TODO: make this concurrent, but its already pretty fast
	for id, b := range boards {
		if existing, ok := storeByID[id]; ok {
			if boardChanged(existing, b) {
				slog.Debug("board sync: update needed for board", "id", id, "name", b.Name)
				nb := models.Board{
					ID:   id,
					Name: b.Name,
				}

				if _, err := txSvc.Boards.Upsert(ctx, nb); err != nil {
					return fmt.Errorf("updating board %d: %w", id, err)
				}
				slog.Debug("board sync: board updated", "id", id, "name", b.Name)
			}
		} else {
			nb := models.Board{
				ID:   id,
				Name: b.Name,
			}

			if _, err := txSvc.Boards.Upsert(ctx, nb); err != nil {
				return fmt.Errorf("board sync: inserting board %d: %w", id, err)
			}
			slog.Debug("board sync: board inserted", "id", id, "name", b.Name)
		}
	}

	for _, d := range sb {
		if _, ok := boards[d.ID]; !ok {
			if err := txSvc.Boards.Delete(ctx, d.ID); err != nil {
				return fmt.Errorf("board sync: deleting board %d: %w", d.ID, err)
			}
			slog.Debug("board sync: deleted board", "id", d.ID, "name", d.Name)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("committing tx: %w", err)
	}

	slog.Debug("board sync complete", "took_time", time.Since(start))
	return nil
}

func boardIDParam(ids []int) string {
	if len(ids) == 0 {
		return ""
	}

	param := ""
	for i, id := range ids {
		param += fmt.Sprintf("board/id = %d", id)
		if i < len(ids)-1 {
			param += " OR "
		}
	}

	return fmt.Sprintf("(%s)", param)
}

func boardChanged(existing models.Board, cwBoard psa.Board) bool {
	return existing.Name != cwBoard.Name
}
