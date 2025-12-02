package notifier

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/thecoretg/ticketbot/internal/models"
)

func (s *Service) getRecipientEmails(ctx context.Context, ticket *models.FullTicket) []string {
	var excluded []models.Member

	// if the sender of the note is a member, exclude them from messages;
	// they don't need a notification for their own note
	if ticket.LatestNote != nil && ticket.LatestNote.Member != nil {
		excluded = append(excluded, *ticket.LatestNote.Member)
	}

	var emails []string
	for _, r := range ticket.Resources {
		if memberSliceContains(excluded, r) {
			continue
		}

		e, err := s.forwardsToEmails(ctx, r.PrimaryEmail)
		if err != nil {
			slog.Error("notifier: error checking forwards for email", "ticket_id", ticket.Ticket.ID, "email", r.PrimaryEmail, "error", err)
		}

		emails = append(emails, e...)
	}

	return filterDuplicateEmails(emails)
}

func (s *Service) forwardsToEmails(ctx context.Context, srcRoom models.WebexRoom) ([]models.WebexRoom, error) {
	rooms, err := s.Rooms.ListByEmail(ctx, email)
	if err != nil {
		return noFwdSlice, fmt.Errorf("listing webex rooms by email: %w", err)
	}

	if len(rooms) == 0 {
		return noFwdSlice, fmt.Errorf("no rooms found by email %s", email)
	}

	var room models.WebexRoom
	if len(rooms) > 1 {
		sort.Slice(rooms, func(i, j int) bool {
			return rooms[i].LastActivity.After(rooms[j].LastActivity)
		})
	}
	room = rooms[0]
	noFwdSlice := []string{email}

	fwds, err := s.Forwards.ListBySourceRoomID(ctx, room.ID)
	if err != nil {
		return noFwdSlice, fmt.Errorf("checking forwards: %w", err)
	}

	if len(fwds) == 0 {
		return noFwdSlice, nil
	}

	activeFwds := filterActiveFwds(fwds)
	if len(activeFwds) == 0 {
		return noFwdSlice, nil
	}

	var emails []string
	for _, f := range activeFwds {
		if f.UserKeepsCopy {
			emails = append(emails, email)
			break
		}
	}

	for _, f := range activeFwds {
		emails = append(emails, f.DestEmail)
	}

	return emails, nil
}

func memberSliceContains(members []models.Member, check models.Member) bool {
	for _, x := range members {
		if x.ID == check.ID {
			return true
		}
	}

	return false
}

func filterActiveFwds(fwds []models.UserForward) []models.UserForward {
	var activeFwds []models.UserForward
	for _, f := range fwds {
		if f.Enabled && dateRangeActive(f.StartDate, f.EndDate) {
			activeFwds = append(activeFwds, f)
		}
	}

	return activeFwds
}

func dateRangeActive(start, end *time.Time) bool {
	now := time.Now()
	if start == nil {
		return false
	}

	if end == nil {
		return now.After(*start)
	}

	return now.After(*start) && now.Before(*end)
}

func filterDuplicateEmails(emails []string) []string {
	seenEmails := make(map[string]struct{})
	for _, e := range emails {
		if _, ok := seenEmails[e]; !ok {
			seenEmails[e] = struct{}{}
		}
	}

	var uniqueEmails []string
	for e := range seenEmails {
		uniqueEmails = append(uniqueEmails, e)
	}

	return uniqueEmails
}
