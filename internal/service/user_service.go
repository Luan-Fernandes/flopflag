package service

import (
	"context"
	"errors"
	"net/http"

	"boliplate-go/internal/domain"
	"boliplate-go/internal/httpx"
	"boliplate-go/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService struct {
	users *repository.UserRepository
}

func NewUserService(users *repository.UserRepository) *UserService {
	return &UserService{users: users}
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpx.NewAppError(http.StatusNotFound, "usuário não encontrado", nil)
		}
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}

func (s *UserService) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	out, err := s.users.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	for i := range out {
		out[i].PasswordHash = ""
	}
	return out, nil
}
