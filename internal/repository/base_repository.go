package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		db: db,
	}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) GetByID(ctx context.Context, id uint64) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error

	if err != nil {
		return nil, err
	}

	return &entity, nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	if entity == nil {
		return errors.New("entity is nil")
	}

	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(new(T), int(id)).Error
}

func (r *BaseRepository[T]) List(ctx context.Context, offset int, limit int) ([]*T, error) {
	var entities []*T

	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&entities).Error

	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(new(T)).Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}
