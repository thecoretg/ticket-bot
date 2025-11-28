package ticketbot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/thecoretg/ticketbot/internal/models"
	"github.com/thecoretg/ticketbot/internal/service/cwsvc"
	"github.com/thecoretg/ticketbot/internal/service/notifier"
)

type Service struct {
	Cfg      *models.Config
	CW       *cwsvc.Service
	Notifier *notifier.Service
}

func New(cfg *models.Config, cw *cwsvc.Service, ns *notifier.Service) *Service {
	return &Service{
		Cfg:      cfg,
		CW:       cw,
		Notifier: ns,
	}
}

func (s *Service) ProcessTicket(ctx context.Context, id int) error {
	exists, err := s.CW.Tickets.Exists(ctx, id)
	if err != nil {
		return fmt.Errorf("checking if ticket exists: %w", err)
	}
	isNew := !exists

	ticket, err := s.CW.ProcessTicket(ctx, id)
	if err != nil {
		return fmt.Errorf("processing ticket: %w", err)
	}

	if s.Cfg.AttemptNotify {
		slog.Debug("ticketbot: attempt notify enabled", "ticket_id", id)
		s.Notifier.ProcessTicket(ctx, ticket, isNew)
		return nil
	}

	slog.Debug("ticketbot: attempt notify disabled", "ticket_id", id)
	return nil
}
