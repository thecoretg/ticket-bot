package ticket

import (
	"context"
	"errors"
	"fmt"
	"time"

	board "github.com/thecoretg/ticketbot/internal/repo/cw_board"
	company "github.com/thecoretg/ticketbot/internal/repo/cw_company"
	contact "github.com/thecoretg/ticketbot/internal/repo/cw_contact"
	member "github.com/thecoretg/ticketbot/internal/repo/cw_member"
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
