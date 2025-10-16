package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/thecoretg/ticketbot/internal/db"
)

func (s *Server) ensureMemberByIdentifier(ctx context.Context, identifier string) (db.CwMember, error) {
	member, err := s.Queries.GetMemberByIdentifier(ctx, identifier)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Debug("member not in store, attempting insert", "member_identifier", identifier)
			cwMember, err := s.CWClient.GetMemberByIdentifier(identifier)
			if err != nil {
				return db.CwMember{}, fmt.Errorf("getting member from cw by identifier: %w", err)
			}

			if cwMember == nil {
				return db.CwMember{}, fmt.Errorf("member %s not found", identifier)
			}

			return s.ensureMemberInStore(ctx, cwMember.ID)
		}
		return db.CwMember{}, fmt.Errorf("querying db for member: %w", err)
	}

	return member, nil
}

func (s *Server) ensureMemberInStore(ctx context.Context, id int) (db.CwMember, error) {
	member, err := s.Queries.GetMember(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Debug("member not in store, attempting insert", "member_id", id)
			cwMember, err := s.CWClient.GetMember(id, nil)
			if err != nil {
				return db.CwMember{}, fmt.Errorf("getting member from cw: %w", err)
			}
			p := db.InsertMemberParams{
				ID:           id,
				Identifier:   cwMember.Identifier,
				FirstName:    cwMember.FirstName,
				LastName:     cwMember.LastName,
				PrimaryEmail: cwMember.PrimaryEmail,
			}
			slog.Debug("created insert member params", "id", p.ID, "identifier", p.Identifier, "first_name", p.FirstName, "last_name", p.LastName, "primary_email", p.PrimaryEmail)

			member, err = s.Queries.InsertMember(ctx, p)
			if err != nil {
				return db.CwMember{}, fmt.Errorf("inserting member into db: %w", err)
			}
			slog.Debug("inserted member into store", "member_id", member.ID, "member_identifier", member.Identifier)
			return member, nil
		} else {
			return db.CwMember{}, fmt.Errorf("getting member from db: %w", err)
		}
	}

	slog.Debug("got existing member from store", "member_id", member.ID, "member_identifier", member.Identifier)
	return member, nil
}
