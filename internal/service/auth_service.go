package service

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Luan-Fernandes/flopflag/internal/domain"
	"github.com/Luan-Fernandes/flopflag/internal/httpx"
	"github.com/Luan-Fernandes/flopflag/internal/repository"
	"github.com/Luan-Fernandes/flopflag/internal/security"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type AuthService struct {
	users  *repository.UserRepository
	tokens *security.TokenManager
}

func NewAuthService(users *repository.UserRepository, tokens *security.TokenManager) *AuthService {
	return &AuthService{
		users:  users,
		tokens: tokens,
	}
}

type AuthResult struct {
	Token string
	User  domain.User
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return nil, httpx.NewAppError(http.StatusBadRequest, "email e senha são obrigatórios", nil)
	}

	hash, err := security.HashPassword(password)
	if err != nil {
		return nil, err
	}

	count, err := s.users.Count(ctx)
	if err != nil {
		return nil, err
	}

	role := domain.RoleUser
	if count == 0 {
		role = domain.RoleAdmin
	}

	u := domain.User{
		Email:        email,
		PasswordHash: hash,
		Role:         role,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := s.users.Create(ctx, &u); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, httpx.NewAppError(http.StatusConflict, "email já cadastrado", nil)
		}
		return nil, err
	}

	token, err := s.tokens.Sign(u.ID, u.Email, u.Role)
	if err != nil {
		return nil, err
	}

	u.PasswordHash = ""
	return &AuthResult{Token: token, User: u}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return nil, httpx.NewAppError(http.StatusBadRequest, "email e senha são obrigatórios", nil)
	}

	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpx.NewAppError(http.StatusUnauthorized, "credenciais inválidas", nil)
		}
		return nil, err
	}

	if err := security.ComparePasswordHash(password, u.PasswordHash); err != nil {
		return nil, httpx.NewAppError(http.StatusUnauthorized, "credenciais inválidas", nil)
	}

	token, err := s.tokens.Sign(u.ID, u.Email, u.Role)
	if err != nil {
		return nil, err
	}

	u.PasswordHash = ""
	return &AuthResult{Token: token, User: *u}, nil
}
