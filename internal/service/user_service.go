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

type UserService interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error)
	GetUserByID(ctx context.Context, id uint64) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	UpdateUser(ctx context.Context, id uint64, req *dto.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id uint64) error
	ListUsers(ctx context.Context, page int, pageSize int) ([]*model.User, error)
}

type userService struct {
	userRepository  repository.UserRepository
	habitRepository repository.HabitRepository
}

func NewUserService(
	userRepository repository.UserRepository,
	habitRepository repository.HabitRepository,
) UserService {
	return &userService{
		userRepository:  userRepository,
		habitRepository: habitRepository,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error) {
	user, err := s.GetUserByUsername(ctx, req.Username)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Failed to query user", zap.Error(err), zap.String("username", req.Username))
		return nil, err
	}

	if user != nil {
		return nil, errors.New("user already exists")
	}

	user = &model.User{
		Username: req.Username,
		Password: req.Password,
		Status:   0,
		Habits:   nil,
	}

	err = s.userRepository.Create(ctx, user)

	if err != nil {
		return nil, err
	}

	logger.Info("User created successfully",
		zap.String("username", user.Username),
		zap.Uint64("user_id", user.ID),
	)

	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uint64) (*model.User, error) {
	user, err := s.userRepository.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := s.userRepository.GetByUsername(ctx, username)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uint64, req *dto.UpdateUserRequest) (*model.User, error) {
	user, err := s.GetUserByID(ctx, id)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Failed to query user", zap.Error(err), zap.Uint64("id", id))
		return nil, err
	}

	if req == nil {
		return nil, errors.New("request is empty")
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	err = s.userRepository.Update(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint64) error {
	user, err := s.userRepository.GetByID(ctx, id)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Failed to query user", zap.Error(err), zap.Uint64("id", id))
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	habits := user.Habits

	if habits != nil {
		for _, habit := range habits {
			err := s.habitRepository.Delete(ctx, habit.ID)

			if err != nil {
				return err
			}
		}
	}

	err = s.userRepository.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *userService) ListUsers(ctx context.Context, page int, pageSize int) ([]*model.User, error) {
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

	users, err := s.userRepository.List(ctx, offset, pageSize)

	if err != nil {
		return nil, err
	}

	return users, nil
}
