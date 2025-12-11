package respository

import (
	"context"
	"w2learn/pkg/def"
)

var _ HealthRepository = (*healthRepository)(nil)

type HealthRepository interface {
	GetDatabaseStatus(ctx context.Context) (int, error)
}

type healthRepository struct {
}

func NewHealthRepository() HealthRepository {
	return &healthRepository{}
}

func (h *healthRepository) GetDatabaseStatus(ctx context.Context) (int, error) {
	return def.HealthStatusCheckError, nil
}
