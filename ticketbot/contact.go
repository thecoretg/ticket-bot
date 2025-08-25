package ticketbot

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/thecoretg/ticketbot/connectwise"
	"github.com/thecoretg/ticketbot/db"
	"log/slog"
)

func (s *Server) ensureContactInStore(ctx context.Context, cwContact *connectwise.Contact) (db.CwContact, error) {
	contact, err := s.Queries.GetContact(ctx, cwContact.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Debug("contact not in store, attempting insert", "contact_id", cwContact.ID, "first_name", cwContact.FirstName, "last_name", cwContact.LastName)
			p := db.InsertContactParams{
				ID:        cwContact.ID,
				FirstName: cwContact.FirstName,
				LastName:  strToPtr(cwContact.LastName),
				CompanyID: intToPtr(cwContact.Company.ID),
			}

			contact, err = s.Queries.InsertContact(ctx, p)
			if err != nil {
				return db.CwContact{}, fmt.Errorf("inserting contact into db: %w", err)
			}
			slog.Info("inserted contact into store", "contact_id", contact.ID, "first_name", contact.FirstName, "last_name", contact.LastName)
			return contact, nil
		} else {
			return db.CwContact{}, fmt.Errorf("getting contact from db: %w", err)
		}
	}

	slog.Debug("got existing contact from store", "contact_id", contact.ID, "first_name", contact.FirstName, "last_name", contact.LastName)
	return contact, nil
}
