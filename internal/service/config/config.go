package config

import (
	"context"

	"github.com/thecoretg/ticketbot/internal/models"
)

type Service struct {
	Config models.ConfigRepository
}

func New(c models.ConfigRepository) *Service {
	return &Service{
		Config: c,
	}
}

func (s *Service) Get(ctx context.Context) (*models.Config, error) {
	return s.Config.Get(ctx)
}

func (s *Service) Update(ctx context.Context, p *models.Config) (*models.Config, error) {
	return s.Config.Upsert(ctx, p)
}
