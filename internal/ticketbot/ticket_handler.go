package ticketbot

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"tctg-automation/pkg/connectwise"
	"time"

	"github.com/gin-gonic/gin"
)

type cwData struct {
	ticket *connectwise.Ticket
	note   *connectwise.ServiceTicketNote
}

func (s *server) addHooksGroup() {
	hooks := s.ginEngine.Group("/hooks")
	cw := hooks.Group("/cw", requireValidCWSignature(), ErrorHandler(s.config.ExitOnError))
	cw.POST("/tickets", s.handleTickets)
}

func (s *server) handleTickets(c *gin.Context) {
	w := &connectwise.WebhookPayload{}
	if err := c.ShouldBindJSON(w); err != nil {
		c.Error(fmt.Errorf("unmarshaling connectwise webhook payload: %w", err))
		return
	}

	if w.ID == 0 {
		c.Error(errors.New("ticket ID cannot be 0"))
		return
	}

	slog.Info("received ticket webhook", "id", w.ID, "action", w.Action)
	if w.Action == "added" || w.Action == "updated" {
		if err := s.addOrUpdateTicket(c.Request.Context(), w.ID, w.Action, false); err != nil {
			c.Error(fmt.Errorf("adding or updating the ticket into data storage: %w", err))
			return
		}

		c.Status(http.StatusNoContent)
	} else {
		c.Status(http.StatusNoContent)
	}
}

func (s *server) getTicketLock(ticketID int) *sync.Mutex {
	lockIface, _ := s.ticketLocks.LoadOrStore(ticketID, &sync.Mutex{})
	return lockIface.(*sync.Mutex)
}

// addOrUpdateTicket serves as the primary handler for updating the data store with ticket data. It also will handle
// extra functionality such as ticket notifications.
func (s *server) addOrUpdateTicket(ctx context.Context, ticketID int, action string, assumeNotify bool) error {
	// Lock the ticket so that extra calls don't interfere. Due to the nature of Connectwise updates will often
	// result in other hooks and actions taking place, which means a ticket rarely only sends one webhook payload.
	lock := s.getTicketLock(ticketID)
	if !lock.TryLock() {
		lock.Lock()
	}

	defer func() {
		lock.Unlock()
	}()

	// Get existing ticket from store - will be nil if it doesn't already exist.
	storeTicket, err := s.dataStore.GetTicket(ticketID)
	if err != nil {
		return fmt.Errorf("getting ticket from storage: %w", err)
	}

	// Get the current data for the ticket via the Connectwise API.
	// This will be used to compare for changes with the store ticket.
	cwData, err := s.getCwData(ticketID)
	if err != nil {
		return fmt.Errorf("getting data from connectwise: %w", err)
	}

	// Get the board the ticket's in from the store - will be nil if it doesn't already exist.
	board, err := s.dataStore.GetBoard(cwData.ticket.Board.ID)
	if err != nil {
		return fmt.Errorf("getting board from storage: %w", err)
	}

	// If the board is nil, add it to the store.
	if board == nil {
		board, err = s.addBoard(cwData.ticket.Board.ID)
		if err != nil {
			return err
		}
		slog.Info("added board to store", "board_id", board.ID, "name", board.Name)
	}

	// Convert the ticket data from Connectwise into a store-compatible ticket.
	// If the store ticket is nil, we'll add the current time as the time it was added.
	workingTicket := cwTicketToStoreTicket(cwData.ticket, cwData.note)
	if storeTicket != nil {
		workingTicket.AddedToStore = storeTicket.AddedToStore
	} else {
		workingTicket.AddedToStore = time.Now()
	}

	// Compare the store ticket and the working ticket to see if there are differences.
	// Also check if the most recent note counts as new for notifier purposes.
	ticketChanged, changeList := findChanges(storeTicket, workingTicket)

	// Insert or update the ticket into the store if it didn't exist or if there were changes.
	if ticketChanged {
		if err := s.dataStore.UpsertTicket(workingTicket); err != nil {
			return fmt.Errorf("upserting ticket to store: %w", err)
		}
	}

	note := &TicketNote{}
	if cwData.note.ID != 0 {
		note, err = s.ensureNoteInStore(cwData.note, assumeNotify)
		if err != nil {
			return fmt.Errorf("ensuring note in store: %w", err)
		}
	}

	// Use the action from the CW hook, whether the note is considered new, and if the board
	// has notifications enabled to determine what type of notification will be sent, if any.
	if meetsMessageCriteria(action, note.Notified, board) {
		if err := s.makeAndSendWebexMsgs(ctx, action, workingTicket, cwData.ticket, board, cwData.note); err != nil {
			return fmt.Errorf("processing webex messages: %w", err)
		}
	}

	// Log the result
	if s.config.Debug {
		slog.Debug("ticket processed", "ticket_id", workingTicket.ID, "action", action, "changes_found", changeList, "latest_note_id", note.ID, "assume_notify", assumeNotify,
			"notified", note.Notified, "board_notify_enbabled", board.NotifyEnabled, "meets_message_criteria", meetsMessageCriteria(action, note.Notified, board),
		)
	} else {
		slog.Info("ticket processed", "ticket_id", workingTicket.ID, "action", action, "notified", note.Notified)
	}

	// Always set notified to true if there is a note
	if note.ID != 0 {
		if err := s.setNotified(note, true); err != nil {
			return fmt.Errorf("setting notified to true: %w", err)
		}
	}

	return nil
}

func (s *server) getCwData(ticketID int) (*cwData, error) {
	ticket, err := s.cwClient.GetTicket(ticketID, nil)
	if err != nil {
		return nil, fmt.Errorf("getting ticket: %w", err)
	}

	note, err := s.cwClient.GetMostRecentTicketNote(ticketID)
	if err != nil {
		return nil, fmt.Errorf("getting most recent note: %w", err)
	}

	if note == nil {
		note = &connectwise.ServiceTicketNote{}
	}

	return &cwData{
		ticket: ticket,
		note:   note,
	}, nil
}

func (s *server) getLatestNote(ticketID int) (*connectwise.ServiceTicketNote, error) {
	note, err := s.cwClient.GetMostRecentTicketNote(ticketID)
	if err != nil {
		return nil, fmt.Errorf("getting most recent note from connectwise: %w", err)
	}

	if note == nil {
		note = &connectwise.ServiceTicketNote{}
	}

	return note, nil
}

func (s *server) ensureNoteInStore(cwNote *connectwise.ServiceTicketNote, assumeNotified bool) (*TicketNote, error) {
	note, err := s.dataStore.GetTicketNote(cwNote.ID)
	if err != nil {
		return nil, fmt.Errorf("getting note from store: %w", err)
	}

	if note == nil {
		note, err = s.addNote(cwNote.TicketId, cwNote.ID, assumeNotified)
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

// addBoard adds Connecwise boards to the data store, with a default of
// notifications not enabled.
func (s *server) addBoard(boardID int) (*Board, error) {
	cwBoard, err := s.cwClient.GetBoard(boardID, nil)
	if err != nil {
		return nil, fmt.Errorf("getting board from connectwise: %w", err)
	}

	storeBoard := &Board{
		ID:            cwBoard.ID,
		Name:          cwBoard.Name,
		NotifyEnabled: false,
	}

	if err := s.dataStore.UpsertBoard(storeBoard); err != nil {
		return nil, fmt.Errorf("adding board to store: %w", err)
	}

	return storeBoard, nil
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

// cwTicketToStoreTicket takes a Connectwise ticket info API response and converts it to a
// struct compatible with our data store.
func cwTicketToStoreTicket(cwTicket *connectwise.Ticket, latestNote *connectwise.ServiceTicketNote) *Ticket {
	return &Ticket{
		ID:        cwTicket.ID,
		Summary:   cwTicket.Summary,
		BoardID:   cwTicket.Board.ID,
		OwnerID:   cwTicket.Owner.ID,
		Resources: cwTicket.Resources,
		UpdatedBy: cwTicket.Info.UpdatedBy,
	}
}

// findChanges compares fields in two data store tickets. It returns a bool for if changes were detected,
// and a comma-separated string of the changes it found (or "none")
func findChanges(a, b *Ticket) (bool, string) {
	if a == nil || b == nil {
		return a != b, "one of the tickets is nil"
	}

	var changedValues []string
	if a.ID != b.ID {
		changedValues = append(changedValues, "ID")
	}

	if a.Summary != b.Summary {
		changedValues = append(changedValues, "Summary")
	}

	if a.BoardID != b.BoardID {
		changedValues = append(changedValues, "BoardID")
	}

	if a.OwnerID != b.OwnerID {
		changedValues = append(changedValues, "OwnerID")
	}

	if a.UpdatedBy != b.UpdatedBy {
		changedValues = append(changedValues, "UpdatedBy")
	}

	if a.Resources != b.Resources {
		changedValues = append(changedValues, "Resources")
	}

	changeStr := "none"
	if len(changedValues) > 0 {
		changeStr = strings.Join(changedValues, ", ")
	}
	return len(changedValues) > 0, changeStr
}

// meetsMessageCriteria checks if a message would be allowed to send a notification,
// depending on if it was added or updated, if the note changed, and the board's notification settings.
func meetsMessageCriteria(action string, noteNotified bool, board *Board) bool {
	if action == "added" {
		return board.NotifyEnabled
	}

	if action == "updated" {
		return !noteNotified && board.NotifyEnabled
	}

	return false
}
