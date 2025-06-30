package ticketbot

import (
	"context"
	"fmt"
	"log/slog"
	"tctg-automation/pkg/connectwise"
)

// Logic for the first run of the DB, which loads all companies, contacts, members, and boards into the DB.

func (s *server) loadInitialData(ctx context.Context) error {
	if err := s.loadInitialBoards(ctx); err != nil {
		return fmt.Errorf("loading initial boards: %w", err)
	}

	if err := s.loadInitialMembers(ctx); err != nil {
		return fmt.Errorf("loading initial members: %w", err)
	}

	if err := s.loadInitialCompanies(ctx); err != nil {
		return fmt.Errorf("loading initial companies: %w", err)
	}

	//if err := s.loadInitialContacts(ctx); err != nil {
	//	return fmt.Errorf("loading initial contacts: %w", err)
	//}

	return nil
}

func (s *server) loadInitialBoards(ctx context.Context) error {
	boards, err := s.cwClient.ListBoards(ctx, nil)
	if err != nil {
		return fmt.Errorf("listing boards from connectwise: %w", err)
	}

	if len(boards) == 0 {
		return nil
	}

	slog.Info("got boards from connectwise", "total", len(boards))
	var dbBoards []Board
	for _, b := range boards {
		d := NewBoard(b.ID, b.Name)
		dbBoards = append(dbBoards, *d)
	}

	tx, err := s.dbHandler.db.Beginx()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	stmt, err := tx.PrepareNamed(upsertBoardSQL())
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("preparing upsert statement: %w", err)
	}
	defer stmt.Close()

	for _, b := range dbBoards {
		if _, err := stmt.Exec(b); err != nil {
			tx.Rollback()
			return fmt.Errorf("upserting board %d: %w", b.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	slog.Info("successfully updated boards in db")
	return nil
}

func (s *server) loadInitialMembers(ctx context.Context) error {
	members, err := s.cwClient.ListMembers(ctx, nil)
	if err != nil {
		return fmt.Errorf("listing boards from connectwise: %w", err)
	}

	if len(members) == 0 {
		slog.Warn("got no members")
		return nil
	}
	slog.Info("got members from connectwise", "total", len(members))

	var dbMembers []Member
	for _, b := range members {
		d := NewMember(b.ID, b.Identifier, b.FirstName, b.LastName, b.PrimaryEmail, b.DefaultPhone)
		dbMembers = append(dbMembers, *d)
	}

	tx, err := s.dbHandler.db.Beginx()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	stmt, err := tx.PrepareNamed(upsertMemberSQL())
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("preparing upsert statement: %w", err)
	}
	defer stmt.Close()

	for _, b := range dbMembers {
		if _, err := stmt.Exec(b); err != nil {
			tx.Rollback()
			return fmt.Errorf("upserting member %d: %w", b.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	slog.Info("successfully updated members in db")
	return nil
}

func (s *server) loadInitialCompanies(ctx context.Context) error {
	p := &connectwise.QueryParams{
		Conditions: "status/name = 'Active'",
	}
	companies, err := s.cwClient.ListCompanies(ctx, p)
	if err != nil {
		return fmt.Errorf("listing companies from connectwise: %w", err)
	}

	if len(companies) == 0 {
		slog.Warn("got no companies")
		return nil
	}
	slog.Info("got companies from connectwise", "total", len(companies))

	var dbCompanies []Company
	for _, c := range companies {
		d := NewCompany(c.Id, c.Name)
		dbCompanies = append(dbCompanies, *d)
	}

	tx, err := s.dbHandler.db.Beginx()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}

	stmt, err := tx.PrepareNamed(upsertCompanySQL())
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("preparing upsert statement: %w", err)
	}
	defer stmt.Close()

	for _, c := range dbCompanies {
		if _, err := stmt.Exec(c); err != nil {
			tx.Rollback()
			return fmt.Errorf("upserting company %d: %w", c.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	slog.Info("successfully updated companies in db")
	return nil
}

func (s *server) loadInitialContacts(ctx context.Context) error {
	comps, err := s.dbHandler.ListCompanies()
	if err != nil {
		return fmt.Errorf("getting companies from db: %w", err)
	}
	slog.Info("got companies from db", "total", len(comps))

	for _, c := range comps {
		p := &connectwise.QueryParams{Conditions: fmt.Sprintf("company/id = %d", c.ID)}
		contacts, err := s.cwClient.ListContacts(ctx, p)

		var dbContacts []Contact
		for _, c := range contacts {
			d := NewContact(c.ID, c.FirstName, c.LastName, c.Company.ID)
			dbContacts = append(dbContacts, *d)
		}

		tx, err := s.dbHandler.db.Beginx()
		if err != nil {
			return fmt.Errorf("beginning transaction: %w", err)
		}

		stmt, err := tx.PrepareNamed(upsertContactSQL())
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("preparing upsert statement: %w", err)
		}
		defer stmt.Close()

		for _, c := range dbContacts {
			if _, err := stmt.Exec(c); err != nil {
				tx.Rollback()
				return fmt.Errorf("upserting contact %d: %w", c.ID, err)
			}
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("committing for company %d: %w", c.ID, err)
		}
		slog.Info("successfully updated contacts", "company", c.Name, "total", len(contacts))
	}

	return nil
}
