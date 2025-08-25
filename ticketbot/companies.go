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

func (s *Server) ensureCompanyInStore(ctx context.Context, cwComp *connectwise.Company) (db.CwCompany, error) {
	company, err := s.Queries.GetCompany(ctx, cwComp.Id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Debug("company not in store, attempting insert", "company_id", cwComp.Id, "company_name", cwComp.Name)
			p := db.InsertCompanyParams{
				ID:   cwComp.Id,
				Name: cwComp.Name,
			}

			company, err = s.Queries.InsertCompany(ctx, p)
			if err != nil {
				return db.CwCompany{}, fmt.Errorf("inserting company into db: %w", err)
			}
			slog.Info("inserted company into store", "company_id", company.ID, "company_name", company.Name)
			return company, nil
		} else {
			return db.CwCompany{}, fmt.Errorf("getting company from db: %w", err)
		}
	}

	slog.Debug("got existing company from store", "company_id", company.ID, "company_name", company.Name)
	return company, nil
}
