package repository

import (
	"context"

	"boliplate-go/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

func (r *TodoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	return r.db.WithContext(ctx).Create(todo).Error
}

func (r *TodoRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Todo, error) {
	var out []domain.Todo
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *TodoRepository) FindByIDForUser(ctx context.Context, id, userID uuid.UUID) (*domain.Todo, error) {
	var t domain.Todo
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TodoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	return r.db.WithContext(ctx).Save(todo).Error
}

func (r *TodoRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Todo{}).Error
}
