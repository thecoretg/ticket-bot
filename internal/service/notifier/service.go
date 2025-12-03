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

type Params struct {
	WebexSvc      *webexsvc.Service
	Recipients    models.WebexRecipientRepository
	Notifiers     models.NotifierRuleRepository
	Notifications models.TicketNotificationRepository
	Forwards      models.NotifierForwardRepository
}

func New(cfg *models.Config, p Params, ms models.MessageSender, cwCompanyID string, maxLen int) *Service {
	return &Service{
		Cfg:              cfg,
		WebexSvc:         p.WebexSvc,
		NotifierRules:    p.Notifiers,
		Notifications:    p.Notifications,
		Forwards:         p.Forwards,
		MessageSender:    ms,
		CWCompanyID:      cwCompanyID,
		MaxMessageLength: maxLen,
	}
}
