package ticketbot

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/thecoretg/ticketbot/connectwise"
	"github.com/thecoretg/ticketbot/db"
	"log/slog"
)

func (s *Server) getLatestNoteFromCW(ticketID int) (*connectwise.ServiceTicketNote, error) {
	note, err := s.CWClient.GetMostRecentTicketNote(ticketID)
	if err != nil {
		return nil, fmt.Errorf("getting most recent note from connectwise: %w", err)
	}

	if note == nil {
		slog.Debug("no most recent note found", "ticket_id", ticketID)
		note = &connectwise.ServiceTicketNote{}
	}

	return note, nil
}

func (s *Server) ensureNoteInStore(ctx context.Context, cwData *cwData, overrideNotify bool) (db.CwTicketNote, error) {
	note, err := s.Queries.GetTicketNote(ctx, cwData.note.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Debug("note not found in store, attempting insert", "ticket_id", cwData.ticket.ID, "note_id", cwData.note.ID)
			note, err = s.Queries.InsertTicketNote(ctx, db.InsertTicketNoteParams{
				ID:        cwData.note.ID,
				TicketID:  cwData.note.TicketId,
				Notified:  overrideNotify,
				MemberID:  getMemberID(cwData),
				ContactID: getContactID(cwData),
			})

			if err != nil {
				return db.CwTicketNote{}, fmt.Errorf("inserting ticket note into db: %w", err)
			}

			slog.Info("inserted note into store", "ticket_id", cwData.ticket.ID, "note_id")
			return note, nil

		} else {
			return db.CwTicketNote{}, fmt.Errorf("getting note from store: %w", err)
		}
	}

	slog.Debug("note already in store", "ticket_id", cwData.ticket.ID, "note_id", cwData.note.ID)
	return note, nil
}

func (s *Server) setNotified(ctx context.Context, noteID int, notified bool) error {
	_, err := s.Queries.SetNoteNotified(ctx, db.SetNoteNotifiedParams{
		ID:       noteID,
		Notified: notified,
	})

	if err != nil {
		return err
	}

	return nil
}

func getMemberID(cwData *cwData) *int {
	if cwData.note.Member.ID != 0 {
		return &cwData.note.Member.ID
	}

	return nil
}

func getContactID(cwData *cwData) *int {
	if cwData.note.Contact.ID != 0 {
		return &cwData.note.Contact.ID
	}

	return nil
}
