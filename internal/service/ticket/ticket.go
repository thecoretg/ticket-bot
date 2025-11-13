package ticket

import (
	"context"
	"errors"
	"fmt"

	board "github.com/thecoretg/ticketbot/internal/repo/cw_board"
	company "github.com/thecoretg/ticketbot/internal/repo/cw_company"
	contact "github.com/thecoretg/ticketbot/internal/repo/cw_contact"
	member "github.com/thecoretg/ticketbot/internal/repo/cw_member"
	ticket "github.com/thecoretg/ticketbot/internal/repo/cw_ticket"
	note "github.com/thecoretg/ticketbot/internal/repo/cw_ticket_note"
)

func (s *Service) getCwData(ticketID int) (cwData, error) {
	t, err := s.cwClient.GetTicket(ticketID, nil)
	if err != nil {
		return cwData{}, fmt.Errorf("getting ticket: %w", err)
	}

	n, err := s.cwClient.GetMostRecentTicketNote(ticketID)
	if err != nil {
		return cwData{}, fmt.Errorf("getting most recent ticket note: %w", err)
	}

	return cwData{ticket: t, note: n}, nil
}

func (s *Service) ensureBoard(ctx context.Context, id int) (board.Board, error) {
	b, err := s.Boards.Get(ctx, id)
	if err != nil && !errors.Is(err, board.ErrNotFound) {
		cw, err := s.cwClient.GetBoard(id, nil)
		if err != nil {
			return board.Board{}, fmt.Errorf("getting board from cw: %w", err)
		}

		b, err = s.Boards.Upsert(ctx, board.Board{
			ID:   cw.ID,
			Name: cw.Name,
		})

		if err != nil {
			return board.Board{}, fmt.Errorf("inserting board into store: %w", err)
		}
	}

	return b, nil
}

func (s *Service) ensureCompany(ctx context.Context, id int) (company.Company, error) {
	c, err := s.Companies.Get(ctx, id)
	if err != nil && !errors.Is(err, company.ErrNotFound) {
		cw, err := s.cwClient.GetCompany(id, nil)
		if err != nil {
			return company.Company{}, fmt.Errorf("getting company from cw: %w", err)
		}

		c, err = s.Companies.Upsert(ctx, company.Company{
			ID:   cw.Id,
			Name: cw.Name,
		})

		if err != nil {
			return company.Company{}, fmt.Errorf("inserting company into store: %w", err)
		}
	}

	return c, nil
}

func (s *Service) ensureContact(ctx context.Context, id int) (contact.Contact, error) {
	c, err := s.Contacts.Get(ctx, id)
	if err != nil && !errors.Is(err, contact.ErrNotFound) {
		cw, err := s.cwClient.GetContact(id, nil)
		if err != nil {
			return contact.Contact{}, fmt.Errorf("getting contact from cw: %w", err)
		}

		var compID *int
		if cw.Company.ID != 0 {
			comp, err := s.ensureCompany(ctx, cw.Company.ID)
			if err != nil {
				return contact.Contact{}, fmt.Errorf("ensuring contact's company is in store: %w", err)
			}
			compID = intToPtr(comp.ID)
		}

		c, err = s.Contacts.Upsert(ctx, contact.Contact{
			ID:        cw.ID,
			FirstName: cw.FirstName,
			LastName:  strToPtr(cw.LastName),
			CompanyID: compID,
		})

		if err != nil {
			return contact.Contact{}, fmt.Errorf("inserting contact into store: %w", err)
		}
	}

	return c, nil
}

func (s *Service) ensureMember(ctx context.Context, id int) (member.Member, error) {
	m, err := s.Members.Get(ctx, id)
	if err != nil && !errors.Is(err, member.ErrNotFound) {
		cw, err := s.cwClient.GetMember(id, nil)
		if err != nil {
			return member.Member{}, fmt.Errorf("getting member from cw: %w", err)
		}

		m, err = s.Members.Upsert(ctx, member.Member{
			ID:           cw.ID,
			Identifier:   cw.Identifier,
			FirstName:    cw.FirstName,
			LastName:     cw.LastName,
			PrimaryEmail: cw.PrimaryEmail,
		})

		if err != nil {
			return member.Member{}, fmt.Errorf("inserting member into store: %w", err)
		}
	}

	return m, nil
}

func (s *Service) ensureTicket(ctx context.Context, cd cwData) (ticket.Ticket, error) {
	t, err := s.Tickets.Get(ctx, cd.ticket.ID)
	if err != nil && !errors.Is(err, ticket.ErrNotFound) {
		t, err = s.Tickets.Upsert(ctx, ticket.Ticket{
			ID:        cd.ticket.ID,
			Summary:   cd.ticket.Summary,
			BoardID:   cd.ticket.Board.ID,
			OwnerID:   intToPtr(cd.ticket.Owner.ID),
			CompanyID: cd.ticket.Company.ID,
			ContactID: intToPtr(cd.ticket.Contact.ID),
			Resources: &cd.ticket.Resources,
			UpdatedBy: &cd.ticket.Info.UpdatedBy,
		})

		if err != nil {
			return ticket.Ticket{}, fmt.Errorf("inserting ticket into store: %w", err)
		}
	}

	return t, nil
}

func (s *Service) ensureTicketNote(ctx context.Context, cd cwData) (note.TicketNote, error) {
	memberID, err := s.getMemberID(ctx, cd)
	if err != nil {
		return note.TicketNote{}, fmt.Errorf("getting member data: %w", err)
	}

	contactID, err := s.getContactID(ctx, cd)
	if err != nil {
		return note.TicketNote{}, fmt.Errorf("getting contact data: %w ", err)
	}

	n, err := s.Notes.Get(ctx, cd.note.ID)
	if err != nil && !errors.Is(err, note.ErrNotFound) {
		n, err = s.Notes.Upsert(ctx, note.TicketNote{
			ID:        cd.note.ID,
			TicketID:  cd.note.TicketId,
			MemberID:  memberID,
			ContactID: contactID,
		})

		if err != nil {
			return note.TicketNote{}, fmt.Errorf("inserting note into store: %w", err)
		}
	}

	return n, nil
}

func (s *Service) getContactID(ctx context.Context, cd cwData) (*int, error) {
	if cd.note.Contact.ID != 0 {
		c, err := s.ensureContact(ctx, cd.note.Contact.ID)
		if err != nil {
			return nil, fmt.Errorf("ensuring contact in store: %w", err)
		}

		return intToPtr(c.ID), nil
	}

	return nil, nil
}

func (s *Service) getMemberID(ctx context.Context, cd cwData) (*int, error) {
	if cd.note.Member.ID != 0 {
		c, err := s.ensureMember(ctx, cd.note.Member.ID)
		if err != nil {
			return nil, fmt.Errorf("ensuring member in store: %w", err)
		}

		return intToPtr(c.ID), nil
	}

	return nil, nil
}

func intToPtr(i int) *int {
	if i == 0 {
		return nil
	}
	val := i
	return &val
}

func strToPtr(s string) *string {
	if s == "" {
		return nil
	}
	val := s
	return &val
}
