package repository

import (
	"context"
	"w2learn/internal/model"

	"gorm.io/gorm"
)

var _ HabitRepository = (*habitRepository)(nil)

type HabitRepository interface {
	Create(ctx context.Context, user *model.Habit) error
	GetByID(ctx context.Context, id uint64) (*model.Habit, error)
	Update(ctx context.Context, user *model.Habit) error
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, offset, limit int) ([]*model.Habit, error)
}

type habitRepository struct {
	*BaseRepository[model.Habit]
}

func NewHabitRepository(db *gorm.DB) HabitRepository {
	return &habitRepository{
		BaseRepository: NewBaseRepository[model.Habit](db),
	}
}
