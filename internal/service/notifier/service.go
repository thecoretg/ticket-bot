package notifier

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thecoretg/ticketbot/internal/external/webex"
	"github.com/thecoretg/ticketbot/internal/models"
)

type Service struct {
	Rooms     models.WebexRoomRepository
	Notifiers models.NotifierRepository
	Forwards  models.UserForwardRepository

	pool        *pgxpool.Pool
	webexClient *webex.Client
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

func New(pool *pgxpool.Pool, r models.WebexRoomRepository, n models.NotifierRepository, f models.UserForwardRepository, wc *webex.Client) *Service {
	return &Service{
		Rooms:       r,
		Notifiers:   n,
		Forwards:    f,
		pool:        pool,
		webexClient: wc,
	}
}

func (s *Service) Run(ctx context.Context, ticket *models.FullTicket, newTicket bool) *Result {
	res := newResult()

}
