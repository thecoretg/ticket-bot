package adgsvc

import (
	"github.com/thecoretg/ticketbot/internal/models"
	"github.com/thecoretg/ticketbot/pkg/addigy"
)

type Service struct {
	Alerts       models.AddigyAlertRepository
	AddigyClient *addigy.Client
}
