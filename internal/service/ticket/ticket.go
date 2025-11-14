package ticket

import (
	"context"
	"errors"
	"fmt"

	"github.com/thecoretg/ticketbot/internal/models"
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

func (s *Service) ensureBoard(ctx context.Context, id int) (models.Board, error) {
	b, err := s.Boards.Get(ctx, id)
	if err != nil && !errors.Is(err, models.ErrBoardNotFound) {
		cw, err := s.cwClient.GetBoard(id, nil)
		if err != nil {
			return models.Board{}, fmt.Errorf("getting board from cw: %w", err)
		}

		b, err = s.Boards.Upsert(ctx, models.Board{
			ID:   cw.ID,
			Name: cw.Name,
		})

		if err != nil {
			return models.Board{}, fmt.Errorf("inserting board into store: %w", err)
		}
	}

	return b, nil
}

func (s *Service) ensureCompany(ctx context.Context, id int) (models.Company, error) {
	c, err := s.Companies.Get(ctx, id)
	if err != nil && !errors.Is(err, models.ErrCompanyNotFound) {
		cw, err := s.cwClient.GetCompany(id, nil)
		if err != nil {
			return models.Company{}, fmt.Errorf("getting company from cw: %w", err)
		}

		c, err = s.Companies.Upsert(ctx, models.Company{
			ID:   cw.Id,
			Name: cw.Name,
		})

		if err != nil {
			return models.Company{}, fmt.Errorf("inserting company into store: %w", err)
		}
	}

	return c, nil
}

func (s *Service) ensureContact(ctx context.Context, id int) (models.Contact, error) {
	c, err := s.Contacts.Get(ctx, id)
	if err != nil && !errors.Is(err, models.ErrContactNotFound) {
		cw, err := s.cwClient.GetContact(id, nil)
		if err != nil {
			return models.Contact{}, fmt.Errorf("getting contact from cw: %w", err)
		}

		var compID *int
		if cw.Company.ID != 0 {
			comp, err := s.ensureCompany(ctx, cw.Company.ID)
			if err != nil {
				return models.Contact{}, fmt.Errorf("ensuring contact's company is in store: %w", err)
			}
			compID = intToPtr(comp.ID)
		}

		c, err = s.Contacts.Upsert(ctx, models.Contact{
			ID:        cw.ID,
			FirstName: cw.FirstName,
			LastName:  strToPtr(cw.LastName),
			CompanyID: compID,
		})

		if err != nil {
			return models.Contact{}, fmt.Errorf("inserting contact into store: %w", err)
		}
	}

	return c, nil
}

func (s *Service) ensureMember(ctx context.Context, id int) (models.Member, error) {
	m, err := s.Members.Get(ctx, id)
	if err != nil && !errors.Is(err, models.ErrMemberNotFound) {
		cw, err := s.cwClient.GetMember(id, nil)
		if err != nil {
			return models.Member{}, fmt.Errorf("getting member from cw: %w", err)
		}

		m, err = s.Members.Upsert(ctx, models.Member{
			ID:           cw.ID,
			Identifier:   cw.Identifier,
			FirstName:    cw.FirstName,
			LastName:     cw.LastName,
			PrimaryEmail: cw.PrimaryEmail,
		})

		if err != nil {
			return models.Member{}, fmt.Errorf("inserting member into store: %w", err)
		}
	}

	return m, nil
}

func (s *Service) ensureTicket(ctx context.Context, cd cwData) (models.Ticket, error) {
	t, err := s.Tickets.Get(ctx, cd.ticket.ID)
	if err != nil && !errors.Is(err, models.ErrTicketNotFound) {
		t, err = s.Tickets.Upsert(ctx, models.Ticket{
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
			return models.Ticket{}, fmt.Errorf("inserting ticket into store: %w", err)
		}
	}

	return t, nil
}

func (s *Service) ensureTicketNote(ctx context.Context, cd cwData) (models.TicketNote, error) {
	memberID, err := s.getMemberID(ctx, cd)
	if err != nil {
		return models.TicketNote{}, fmt.Errorf("getting member data: %w", err)
	}

	contactID, err := s.getContactID(ctx, cd)
	if err != nil {
		return models.TicketNote{}, fmt.Errorf("getting contact data: %w ", err)
	}

	n, err := s.Notes.Get(ctx, cd.note.ID)
	if err != nil && !errors.Is(err, models.ErrTicketNoteNotFound) {
		n, err = s.Notes.Upsert(ctx, models.TicketNote{
			ID:        cd.note.ID,
			TicketID:  cd.note.TicketId,
			MemberID:  memberID,
			ContactID: contactID,
		})

		if err != nil {
			return models.TicketNote{}, fmt.Errorf("inserting note into store: %w", err)
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
