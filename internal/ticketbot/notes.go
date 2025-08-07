package ticketbot

import (
	"fmt"
	"tctg-automation/pkg/connectwise"
)

func (s *server) getLatestNoteFromCW(ticketID int) (*connectwise.ServiceTicketNote, error) {
	note, err := s.cwClient.GetMostRecentTicketNote(ticketID)
	if err != nil {
		return nil, fmt.Errorf("getting most recent note from connectwise: %w", err)
	}

	if note == nil {
		note = &connectwise.ServiceTicketNote{}
	}

	return note, nil
}

func (s *server) ensureNoteInStore(cwData *cwData, assumeNotified bool) (*TicketNote, error) {
	note, err := s.dataStore.GetTicketNote(cwData.note.ID)
	if err != nil {
		return nil, fmt.Errorf("getting note from store: %w", err)
	}

	if note == nil {
		note, err = s.addNote(cwData.ticket.ID, cwData.note.ID, assumeNotified)
		if err != nil {
			return nil, fmt.Errorf("adding note to store: %w", err)
		}
	}

	return note, nil
}

func (s *server) setNotified(note *TicketNote, notified bool) error {
	note.Notified = notified
	if err := s.dataStore.UpsertTicketNote(note); err != nil {
		return fmt.Errorf("upserting note: %w", err)
	}

	return nil
}

// addNote adds a ticket note to the data store
func (s *server) addNote(ticketID, noteID int, notified bool) (*TicketNote, error) {
	note := &TicketNote{
		ID:       noteID,
		TicketID: ticketID,
		Notified: notified,
	}

	if err := s.dataStore.UpsertTicketNote(note); err != nil {
		return nil, err
	}

	return note, nil
}
