package ticketbot

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/thecoretg/ticketbot/db"
	"log/slog"
)

func (s *Server) ensureMemberInStore(ctx context.Context, id int, identifier, firstName, lastName, email string) (db.CwMember, error) {
	member, err := s.Queries.GetMember(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Debug("member not in store, attempting insert", "member_id", id)
			p := db.InsertMemberParams{
				ID:           id,
				Identifier:   identifier,
				FirstName:    firstName,
				LastName:     lastName,
				PrimaryEmail: email,
			}

			member, err = s.Queries.InsertMember(ctx, p)
			if err != nil {
				return db.CwMember{}, fmt.Errorf("inserting member into db: %w", err)
			}
			slog.Info("inserted member into store", "member_id", member.ID, "member_identifier", member.Identifier)
			return member, nil
		} else {
			return db.CwMember{}, fmt.Errorf("getting member from db: %w", err)
		}
	}

	slog.Debug("got existing member from store", "member_id", member.ID, "member_identifier", member.Identifier)
	return member, nil
}
