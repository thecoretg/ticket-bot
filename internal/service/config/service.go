package config

import (
	"context"
	"errors"

	"github.com/thecoretg/ticketbot/internal/models"
)

var (
	ErrNoPayload = errors.New("nil payload received")
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
	return s.ensureConfig(ctx)
}

func (s *Service) Update(ctx context.Context, p *models.Config) (*models.Config, error) {
	if p == nil {
		return nil, ErrNoPayload
	}

	return s.Config.Upsert(ctx, p)
}
