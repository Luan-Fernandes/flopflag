package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"boliplate-go/internal/domain"
	"boliplate-go/internal/httpx"
	"boliplate-go/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoService struct {
	todos *repository.TodoRepository
}

func NewTodoService(todos *repository.TodoRepository) *TodoService {
	return &TodoService{todos: todos}
}

func (s *TodoService) List(ctx context.Context, userID uuid.UUID) ([]domain.Todo, error) {
	return s.todos.ListByUserID(ctx, userID)
}

func (s *TodoService) Get(ctx context.Context, id, userID uuid.UUID) (*domain.Todo, error) {
	t, err := s.todos.FindByIDForUser(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpx.NewAppError(http.StatusNotFound, "todo não encontrado", nil)
		}
		return nil, err
	}
	return t, nil
}

func (s *TodoService) Create(ctx context.Context, userID uuid.UUID, title string) (*domain.Todo, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, httpx.NewAppError(http.StatusBadRequest, "title é obrigatório", nil)
	}

	t := domain.Todo{
		UserID: userID,
		Title:  title,
		Done:   false,
	}
	if err := s.todos.Create(ctx, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *TodoService) Update(ctx context.Context, id, userID uuid.UUID, title *string, done *bool) (*domain.Todo, error) {
	t, err := s.todos.FindByIDForUser(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpx.NewAppError(http.StatusNotFound, "todo não encontrado", nil)
		}
		return nil, err
	}

	if title != nil {
		trimmed := strings.TrimSpace(*title)
		if trimmed == "" {
			return nil, httpx.NewAppError(http.StatusBadRequest, "title não pode ser vazio", nil)
		}
		t.Title = trimmed
	}
	if done != nil {
		t.Done = *done
	}

	if err := s.todos.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TodoService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	if _, err := s.todos.FindByIDForUser(ctx, id, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpx.NewAppError(http.StatusNotFound, "todo não encontrado", nil)
		}
		return err
	}

	if err := s.todos.Delete(ctx, id, userID); err != nil {
		return err
	}
	return nil
}
