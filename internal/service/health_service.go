package service

import (
	"context"
	"w2learn/internal/model"
	"w2learn/internal/respository"
	"w2learn/pkg/def"
)

var _ HealthService = (*healthService)(nil)

type HealthService interface {
	GetHealth(ctx context.Context, flag int) (*model.HealthModel, error)
}

type healthService struct {
	healthRepo respository.HealthRepository
}

func NewHealthService(repository respository.HealthRepository) HealthService {
	return &healthService{
		healthRepo: repository,
	}
}

func (s *healthService) GetHealth(ctx context.Context, flag int) (*model.HealthModel, error) {
	healthModel := model.GetDefaultHealthModel()

	if flag&def.HealthStatusRequestFlagDatabase > 0 {
		status, err := s.healthRepo.GetDatabaseStatus(ctx)

		if err != nil {
			return nil, err
		}

		healthModel.DatabaseStatus = status
	}

	return healthModel, nil
}
