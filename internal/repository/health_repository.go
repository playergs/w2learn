package repository

import (
	"context"
	"w2learn/pkg/database"
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
	db := database.GetDB()

	if db == nil {
		return def.HealthStatusCheckError, nil
	}

	baseDB, err := db.DB()

	if err != nil {
		return def.HealthStatusCheckError, err
	}

	err = baseDB.PingContext(ctx)

	if err != nil {
		return def.HealthStatusCheckError, err
	}

	return def.HealthStatusCheckOK, nil
}
