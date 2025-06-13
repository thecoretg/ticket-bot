package ticketbot

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"tctg-automation/pkg/connectwise"
	"tctg-automation/pkg/util"
	"tctg-automation/pkg/webex"
)

const (
	maxNoteLength = 300
)

func (s *server) processTicket(ctx context.Context, w *connectwise.WebhookPayload) (int, map[string]any) {
	slog.Debug("received webhook payload", "action", w.Action, "ticketID", w.ID, "memberID", w.MemberId)

	if !validAction(w.Action) {
		return http.StatusOK, util.ResultJSON(fmt.Sprintf("invalid action: %s", w.Action))
	}

	if w.ID == 0 {
		return http.StatusBadRequest, util.ErrorJSON("ticket id not provided")
	}

	t, err := s.cwClient.GetTicket(ctx, w.ID, nil)
	if err != nil {
		return http.StatusInternalServerError, util.ErrorJSON(fmt.Sprintf("getting ticket details: %v", err))
	}

	b, err := getBoard(s.db, t.Board.ID)
	if err != nil {
		return http.StatusInternalServerError, util.ErrorJSON(fmt.Sprintf("getting board details for ticket: %v", err))
	}

	if b == nil {
		return http.StatusOK, util.ResultJSON(fmt.Sprintf("board %d not found in settings", t.Board.ID))
	}

	if !b.Enabled {
		return http.StatusOK, util.ResultJSON(fmt.Sprintf("board %d is not enabled for notifications", t.Board.ID))
	}

	p := &connectwise.QueryParams{OrderBy: "_info/dateEntered desc"}
	n, err := s.cwClient.ListServiceTicketNotes(ctx, t.ID, p)
	if err != nil {
		return http.StatusInternalServerError, util.ErrorJSON(fmt.Sprintf("getting notes for ticket: %v", err))
	}

	switch w.Action {
	case "added":
		return s.processNewTicket(ctx, b, t, n)
	case "updated":
		return s.processUpdatedTicket(ctx, w, t, n)
	}

	return http.StatusInternalServerError, util.ErrorJSON("not sure how we got past the gate, but we're here")
}

func (s *server) processNewTicket(ctx context.Context, b *boardSetting, t *connectwise.Ticket, n []connectwise.ServiceTicketNoteAll) (int, map[string]any) {
	if b.WebexRoomID == "" {
		return http.StatusInternalServerError, util.ErrorJSON(fmt.Sprintf("board %d enabled, but no webex room ID found", b.BoardID))
	}

	m := buildNewTicketMessage(t, n)
	w := webex.NewMessageToRoom(b.WebexRoomID, m)
	if err := s.webexClient.SendMessage(ctx, w); err != nil {
		return http.StatusInternalServerError, util.ErrorJSON(fmt.Sprintf("sending message to webex room: %v", err))
	}

	return http.StatusNoContent, nil
}

func (s *server) processUpdatedTicket(ctx context.Context, w *connectwise.WebhookPayload, t *connectwise.Ticket, n []connectwise.ServiceTicketNoteAll) (int, map[string]any) {
	var mutedUsers []string

	if t.Resources == "" {
		return http.StatusOK, util.ResultJSON(fmt.Sprintf("ticket %d has no resources", t.ID))
	}

	if updatedBy, present := hasUpdatedBy(w.Entity); present {
		mutedUsers = append(mutedUsers, updatedBy)

		u, err := getUser(s.db, updatedBy)
		if err != nil {
			return http.StatusInternalServerError, util.ErrorJSON(fmt.Sprintf("getting user by connectwise id: %v", err))
		}

		if u != nil && u.IgnoreUpdate {
			return http.StatusOK, util.ResultJSON(fmt.Sprintf("user %s marked as ignore updates, no message to send", updatedBy))
		}
	}

	ln := mostRecentNote(n)
	if ln != nil {
		mutedUsers = append(mutedUsers, noteSenderName(ln))
	}

	e, err := s.getAndStoreResourceEmails(ctx, t.Resources, mutedUsers)
	if err != nil {
		return http.StatusInternalServerError, util.ErrorJSON(fmt.Sprintf("getting emails from resource ids: %v", err))
	}

	if len(e) == 0 {
		return http.StatusOK, util.ResultJSON("no resources to send messages to")
	}

	m := buildUpdatedTicketMessage(t, n)
	if err := s.sendUpdateMessages(ctx, m, e); err != nil {
		return http.StatusInternalServerError, util.ErrorJSON(fmt.Sprintf("sending messages to users: %v", err))
	}

	return http.StatusNoContent, nil
}

func (s *server) getAndStoreResourceEmails(ctx context.Context, resourceString string, mutedUsers []string) ([]string, error) {
	ids := splitTicketResources(resourceString)
	if ids == nil {
		return nil, nil // No resources to process
	}

	var emails []string
	for _, id := range ids {
		if isMuted(id, mutedUsers) {
			continue
		}

		id = strings.TrimSpace(id)
		if id == "" {
			continue // Skip empty IDs
		}

		u, err := getUser(s.db, id)
		if err != nil {
			continue
		}

		if u != nil {
			if u.Email == "" {
				continue // Skip users without an email, no need for mute check
			}

			if u.Mute {
				continue // Skip excluded users
			}

			emails = append(emails, u.Email)
			continue // Use cached email if available
		}

		email, err := s.getMemberEmail(ctx, id)
		if err != nil {
			continue
		}

		if email == "" {
			continue // Skip if no email is found
		}

		newUser := &user{
			CWId:  id,
			Email: email,
			Mute:  false, // Default to not excluded
		}

		if err := createOrUpdateUser(s.db, newUser); err != nil {
			continue
		}

		emails = append(emails, email)
	}

	return emails, nil
}

func (s *server) getMemberEmail(ctx context.Context, id string) (string, error) {
	q := &connectwise.QueryParams{
		Conditions: fmt.Sprintf("Identifier='%s'", id),
	}

	m, err := s.cwClient.ListMembers(ctx, q)
	if err != nil {
		return "", fmt.Errorf("getting member with query: %w", err)
	}

	if len(m) == 0 {
		return "", nil
	}

	if len(m) > 1 {
		return "", fmt.Errorf("too many members (%d) returned for id %s", len(m), id)
	}

	if m[0].PrimaryEmail == "" {
		return "", fmt.Errorf("empty email found for member id %s", id)
	}

	return m[0].PrimaryEmail, nil
}

func (s *server) sendUpdateMessages(ctx context.Context, m string, emails []string) error {
	for _, e := range emails {
		w := webex.NewMessageToPerson(e, m)
		if err := s.webexClient.SendMessage(ctx, w); err != nil {
			return fmt.Errorf("sending webex message to %s: %w", e, err)
		}
	}

	return nil
}

func hasUpdatedBy(entity string) (string, bool) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(entity), &m); err != nil {
		return "", false
	}

	info, ok := m["_info"].(map[string]interface{})
	if !ok {
		return "", false
	}

	val, present := info["updatedBy"]
	if !present {
		return "", false
	}

	strVal, ok := val.(string)
	return strVal, strVal != ""
}

func buildNewTicketMessage(t *connectwise.Ticket, n []connectwise.ServiceTicketNoteAll) string {
	m := fmt.Sprintf("**New:** %s %s", ticketLink(t), t.Summary)

	// add requester contact and company name, or just company if no contact (rare)
	r := fmt.Sprintf("**Requester:** %s (No Contact)", t.Company.Name)
	if t.Contact.Name != "" {
		r = fmt.Sprintf("**Requester:** %s (%s)", t.Contact.Name, t.Company.Name)
	}
	m += fmt.Sprintf("\n%s", r)

	// add most recent note if present
	mr := mostRecentNote(n)
	if mr != nil {
		// trim note text and add ... if it exceeds the maximum
		noteTxt := mr.Text
		if len(noteTxt) > maxNoteLength {
			noteTxt = noteTxt[:maxNoteLength] + "..."
		}

		m += fmt.Sprintf("\n**Latest Note:** %s\n"+
			"%s",
			noteSenderName(mr),
			addBlockQuotes(noteTxt),
		)
	}

	return m
}

func ticketLink(t *connectwise.Ticket) string {
	return fmt.Sprintf("[%d](https://na.myconnectwise.net/v4_6_release/services/system_io/Service/fv_sr100_request.rails?service_recid=%d&companyName=securenetit)", t.ID, t.ID)
}

func buildUpdatedTicketMessage(t *connectwise.Ticket, n []connectwise.ServiceTicketNoteAll) string {
	m := fmt.Sprintf("**New Response:** %s %s", ticketLink(t), t.Summary)

	// add most recent note if present
	mr := mostRecentNote(n)
	if mr != nil {
		// trim note text and add ... if it exceeds the maximum
		noteTxt := mr.Text
		if len(noteTxt) > maxNoteLength {
			noteTxt = noteTxt[:maxNoteLength] + "..."
		}

		m += fmt.Sprintf("\n**Latest Note:** %s\n"+
			"%s",
			noteSenderName(mr),
			addBlockQuotes(noteTxt),
		)
	}

	return m
}

func noteSenderName(n *connectwise.ServiceTicketNoteAll) string {
	if n.Member.Name != "" {
		return n.Member.Name
	} else if n.Contact.Name != "" {
		return n.Contact.Name
	} else {
		return "N/A"
	}
}

func mostRecentNote(n []connectwise.ServiceTicketNoteAll) *connectwise.ServiceTicketNoteAll {
	for _, note := range n {
		if note.Text != "" {
			return &note
		}
	}

	return nil
}

func splitTicketResources(resourceString string) []string {
	parts := strings.Split(resourceString, ",")
	var resources []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		resources = append(resources, part)
	}

	if len(resources) == 0 {
		return nil
	} else {
		return resources
	}
}

func addBlockQuotes(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line == "" {
			lines[i] = ">"
		} else {
			lines[i] = "> " + line
		}
	}

	return strings.Join(lines, "\n")
}

func validAction(action string) bool {
	return action == "added" || action == "updated"
}

func isMuted(email string, mutedUsers []string) bool {
	for _, e := range mutedUsers {
		if e == email {
			return true
		}
	}
	return false
}
