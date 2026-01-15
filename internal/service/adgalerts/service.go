package adgalerts

import (
	"github.com/thecoretg/ticketbot/internal/service/adgsvc"
	"github.com/thecoretg/ticketbot/internal/service/cwsvc"
	"github.com/thecoretg/ticketbot/pkg/addigy"
	"github.com/thecoretg/ticketbot/pkg/psa"
)

type Service struct {
	AddigyClient *addigy.Client
	AddigySvc    *adgsvc.Service
	CWClient     *psa.Client
	CWSvc        *cwsvc.Service
}

func New(client *addigy.Client, addigySvc *adgsvc.Service, cwClient *psa.Client, cwSvc *cwsvc.Service) *Service {
	return &Service{
		AddigyClient: client,
		AddigySvc:    addigySvc,
		CWClient:     cwClient,
		CWSvc:        cwSvc,
	}
}
