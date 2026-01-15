package adgsvc

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/thecoretg/ticketbot/internal/models"
)

type Service struct {
	Alerts models.AddigyAlertRepository
}

func New(r models.AddigyAlertRepository) *Service {
	return &Service{
		Alerts: r,
	}
}

func (s *Service) WithTX(tx pgx.Tx) *Service {
	return &Service{
		Alerts: s.Alerts.WithTx(tx),
	}
}

func (s *Service) ListAlerts(ctx context.Context) ([]*models.AddigyAlert, error) {
	return s.Alerts.List(ctx)
}

func (s *Service) ListAlertsByStatus(ctx context.Context, status string) ([]*models.AddigyAlert, error) {
	return s.Alerts.ListByStatus(ctx, status)
}

func (s *Service) ListUnresolvedAlerts(ctx context.Context) ([]*models.AddigyAlert, error) {
	return s.Alerts.ListUnresolved(ctx)
}

func (s *Service) ListAlertsByTicket(ctx context.Context, ticketID int) ([]*models.AddigyAlert, error) {
	return s.Alerts.ListByTicket(ctx, ticketID)
}

func (s *Service) GetAlert(ctx context.Context, id string) (*models.AddigyAlert, error) {
	return s.Alerts.Get(ctx, id)
}

func (s *Service) CreateAlert(ctx context.Context, a *models.AddigyAlert) (*models.AddigyAlert, error) {
	return s.Alerts.Create(ctx, a)
}

func (s *Service) UpsertAlert(ctx context.Context, a *models.AddigyAlert) (*models.AddigyAlert, error) {
	return s.Alerts.Upsert(ctx, a)
}

func (s *Service) UpdateAlertTicket(ctx context.Context, id string, ticketID *int) error {
	return s.Alerts.UpdateTicket(ctx, id, ticketID)
}

func (s *Service) UpdateAlertStatus(ctx context.Context, id string, status string) error {
	return s.Alerts.UpdateStatus(ctx, id, status)
}

func (s *Service) AcknowledgeAlert(ctx context.Context, id string, acknowledgedOn time.Time) error {
	return s.Alerts.Acknowledge(ctx, id, acknowledgedOn)
}

func (s *Service) ResolveAlert(ctx context.Context, id string, resolvedOn time.Time, resolvedByEmail string) error {
	return s.Alerts.Resolve(ctx, id, resolvedOn, resolvedByEmail)
}

func (s *Service) DeleteAlert(ctx context.Context, id string) error {
	return s.Alerts.Delete(ctx, id)
}
