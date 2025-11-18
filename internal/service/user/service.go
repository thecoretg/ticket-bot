package user

import (
	"context"

	"github.com/thecoretg/ticketbot/internal/models"
)

type Service struct {
	Users models.APIUserRepository
	Keys  models.APIKeyRepository
}

func New(u models.APIUserRepository, k models.APIKeyRepository) *Service {
	return &Service{
		Users: u,
		Keys:  k,
	}
}

func (s *Service) List(ctx context.Context) ([]models.APIUser, error) {
	return s.Users.List(ctx)
}

func (s *Service) Get(ctx context.Context, id int) (*models.APIUser, error) {
	return s.Users.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.Users.Delete(ctx, id)
}

func (s *Service) AddAPIKey(ctx context.Context, userEmail string) (string, error) {
	k, err := s.createAPIKey(ctx, userID)
}
