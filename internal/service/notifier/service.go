package notifier

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/models"
	"github.com/thecoretg/ticketbot/internal/service/webexsvc"
)

type Service struct {
	Cfg              *models.Config
	WebexSvc         *webexsvc.Service
	NotifierRules    models.NotifierRuleRepository
	Notifications    models.TicketNotificationRepository
	Forwards         models.NotifierForwardRepository
	Pool             *pgxpool.Pool
	MessageSender    models.MessageSender
	CWCompanyID      string
	MaxMessageLength int
}

type SvcParams struct {
	Cfg              *models.Config
	WebexSvc         *webexsvc.Service
	NotifierRules    models.NotifierRuleRepository
	Notifications    models.TicketNotificationRepository
	Forwards         models.NotifierForwardRepository
	Pool             *pgxpool.Pool
	MessageSender    models.MessageSender
	CWCompanyID      string
	MaxMessageLength int
}

func New(p SvcParams) *Service {
	return &Service{
		Cfg:              p.Cfg,
		WebexSvc:         p.WebexSvc,
		NotifierRules:    p.NotifierRules,
		Notifications:    p.Notifications,
		Forwards:         p.Forwards,
		Pool:             p.Pool,
		MessageSender:    p.MessageSender,
		CWCompanyID:      p.CWCompanyID,
		MaxMessageLength: p.MaxMessageLength,
	}
}
