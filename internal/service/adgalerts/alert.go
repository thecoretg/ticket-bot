package adgalerts

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/thecoretg/ticketbot/internal/models"
	"github.com/thecoretg/ticketbot/pkg/addigy"
)

const ticketRetryLimit = 5

var (
	ErrNoAlertID         = errors.New("no alert id provided")
	ErrReachedRetryLimit = errors.New("reached ticket id retry limit")
)

func (s *Service) RunInitialTicketUpdates(ctx context.Context, alertID string) error {
	if alertID == "" {
		return ErrNoAlertID
	}
	lg := slog.Default().With(slog.String("alert_id", alertID))
	aa, err := s.getAlertFromAddigy(alertID, lg)
	if err != nil {
		return fmt.Errorf("getting alert from addigy: %w", err)
	}

	lg = lg.With(
		slog.Int("ticket_id", *aa.TicketID),
		slog.String("agent_id", aa.AgentID),
		slog.String("name", aa.Name),
		slog.String("fact_name", aa.FactName),
		slog.Bool("remediation_enabled", aa.RemediationEnabled),
	)
	lg.Debug("got alert from addgy")
	a, err := s.upsertAlertToStore(ctx, aa)

	return nil
}

func (s *Service) getAlertFromAddigy(id string, lg *slog.Logger) (*addigy.Alert, error) {
	if id == "" {
		return nil, ErrNoAlertID
	}

	al, err := s.AddigyClient.GetAlert(id)
	if err != nil {
		return nil, fmt.Errorf("getting alert from addigy: %w", err)
	}

	for i := 1; al.TicketID == nil && i <= ticketRetryLimit; i++ {
		if i == ticketRetryLimit {
			return nil, ErrReachedRetryLimit
		}

		lg.Warn("alert received, but no ticket ID; trying again in 5 seconds", "attempt", i)
		time.Sleep(5 * time.Second)

		al, err = s.AddigyClient.GetAlert(id)
		if err != nil {
			return nil, err
		}
	}

	lg.Debug("got ticket id", "id", al.TicketID)
	return al, nil
}

func (s *Service) upsertAlertToStore(ctx context.Context, aa *addigy.Alert) (*models.AddigyAlert, error) {
	a := &models.AddigyAlert{
		ID:             aa.ID,
		TicketID:       aa.TicketID,
		Level:          aa.Level,
		Category:       aa.Category,
		Name:           aa.Name,
		FactName:       aa.FactName,
		FactIdentifier: aa.FactIdentifier,
		FactType:       aa.ValueType,
		Selector:       aa.Selector,
		Status:         aa.Status,
		Muted:          aa.Muted,
		Remediation:    aa.RemediationEnabled,
		ResolvedOn:     parseStringToTime(aa.ResolvedDate),
		AcknowledgedOn: parseStringToTime(aa.AcknowledgedDate),
	}

	if aa.ResolvedUserEmail != "" {
		a.ResolvedByEmail = &aa.ResolvedUserEmail
	}

	return s.AddigySvc.UpsertAlert(ctx, a)
}

func parseStringToTime(s string) *time.Time {
	if s == "" {
		return nil
	}

	pt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}

	return &pt
}
