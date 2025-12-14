package service

import (
	"context"
	"errors"
	"w2learn/internal/dto"
	"w2learn/internal/model"
	"w2learn/internal/repository"
	"w2learn/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HabitService interface {
	CreateHabit(ctx context.Context, req *dto.CreateHabitRequest) (*model.Habit, error)
	GetHabitByID(ctx context.Context, id uint64) (*model.Habit, error)
	UpdateHabit(ctx context.Context, hid uint64, req *dto.UpdateHabitRequest) (*model.Habit, error)
	DeleteHabit(ctx context.Context, req *dto.DeleteHabitRequest) error
	ListHabits(ctx context.Context, page int, pageSize int) ([]*model.Habit, error)
}

type habitService struct {
	habitRepository repository.HabitRepository
	userRepository  repository.UserRepository
}

func NewHabitService(
	habitRepository repository.HabitRepository,
	userRepository repository.UserRepository,
) HabitService {
	return &habitService{
		habitRepository: habitRepository,
		userRepository:  userRepository,
	}
}

func (s *habitService) CreateHabit(ctx context.Context, req *dto.CreateHabitRequest) (*model.Habit, error) {
	user, err := s.userRepository.GetByID(ctx, req.UserID)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("userRepository.GetByID", zap.Error(err))
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.Habits == nil {
		user.Habits = make([]model.Habit, 0)
	}

	habit := model.Habit{
		UserID: req.UserID,
		Name:   req.Name,
		Info:   req.Info,
	}

	err = s.habitRepository.Create(ctx, &habit)

	if err != nil {
		return nil, err
	}

	user.Habits = append(user.Habits, habit)

	err = s.userRepository.Update(ctx, user)

	if err != nil {
		return nil, err
	}

	return &habit, nil
}

func (s *habitService) GetHabitByID(ctx context.Context, id uint64) (*model.Habit, error) {
	habit, err := s.habitRepository.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return habit, nil
}

func (s *habitService) UpdateHabit(ctx context.Context, hid uint64, req *dto.UpdateHabitRequest) (*model.Habit, error) {
	habit, err := s.habitRepository.GetByID(ctx, hid)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("habitRepository.GetByID", zap.Error(err))
		return nil, err
	}

	if habit == nil {
		return nil, errors.New("habit not found")
	}

	if req == nil {
		return nil, errors.New("req is nil")
	}

	if req.Name != "" {
		habit.Name = req.Name
	}

	if req.Info != "" {
		habit.Info = req.Info
	}

	err = s.habitRepository.Update(ctx, habit)

	if err != nil {
		return nil, err
	}

	return habit, nil
}

func (s *habitService) DeleteHabit(ctx context.Context, req *dto.DeleteHabitRequest) error {
	if req == nil {
		return errors.New("req is nil")
	}

	user, err := s.userRepository.GetByID(ctx, req.UserID)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("userRepository.GetByID", zap.Error(err))
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if user.Habits == nil {
		return errors.New("habits is nil")
	}

	deleteIndex := -1
	var deleteID uint64

	for index, habit := range user.Habits {
		if habit.ID == req.HabitID {
			deleteIndex = index
			deleteID = habit.ID
			break
		}
	}

	if deleteIndex == -1 {
		return errors.New("habit not found")
	}

	if deleteIndex == 0 {
		user.Habits = user.Habits[deleteIndex+1:]
	} else if deleteIndex == len(user.Habits)-1 {
		user.Habits = user.Habits[:deleteIndex]
	} else {
		user.Habits = append(user.Habits[:deleteIndex], user.Habits[deleteIndex+1:]...)
	}

	err = s.userRepository.Update(ctx, user)

	if err != nil {
		return err
	}

	err = s.habitRepository.Delete(ctx, deleteID)

	if err != nil {
		return err
	}

	return nil
}

func (s *habitService) ListHabits(ctx context.Context, page int, pageSize int) ([]*model.Habit, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	list, err := s.habitRepository.List(ctx, offset, pageSize)

	if err != nil {
		return nil, err
	}

	return list, err
}
