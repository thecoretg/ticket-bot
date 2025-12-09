package webexsvc

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/models"
)

type Service struct {
	Recipients  models.WebexRecipientRepository
	pool        *pgxpool.Pool
	WebexClient models.MessageSender
}

func New(pool *pgxpool.Pool, r models.WebexRecipientRepository, cl models.MessageSender) *Service {
	return &Service{
		Recipients:  r,
		WebexClient: cl,
		pool:        pool,
	}
}

func (s *Service) WithTx(tx pgx.Tx) *Service {
	return &Service{
		Recipients:  s.Recipients.WithTx(tx),
		WebexClient: s.WebexClient,
		pool:        s.pool,
	}
}
