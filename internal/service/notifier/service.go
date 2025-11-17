package notifier

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/external/webex"
	"github.com/thecoretg/ticketbot/internal/models"
)

type Service struct {
	Rooms         models.WebexRoomRepository
	Notifiers     models.NotifierRepository
	Notifications models.TicketNotificationRepository
	Forwards      models.UserForwardRepository

	pool             *pgxpool.Pool
	webexClient      *webex.Client
	cwClientID       string
	maxMessageLength int
}

type Result struct {
	MembersToNotify []models.Member
	RoomsToNotify   []models.WebexRoom
	SuccessNotis    []string
	FailedNotis     []string
	Error           error
}

func newResult() *Result {
	return &Result{
		MembersToNotify: []models.Member{},
		RoomsToNotify:   []models.WebexRoom{},
		SuccessNotis:    []string{},
		FailedNotis:     []string{},
		Error:           nil,
	}
}

func New(pool *pgxpool.Pool, r models.WebexRoomRepository, n models.NotifierRepository,
	m models.TicketNotificationRepository, f models.UserForwardRepository, wc *webex.Client, cwClientID string, max int) *Service {
	return &Service{
		Rooms:            r,
		Notifiers:        n,
		Notifications:    m,
		Forwards:         f,
		pool:             pool,
		webexClient:      wc,
		cwClientID:       cwClientID,
		maxMessageLength: max,
	}
}

func (s *Service) Run(ctx context.Context, ticket *models.FullTicket, newTicket bool) *Result {
	res := newResult()

}
