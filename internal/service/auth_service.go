package service

import (
	"context"
	"errors"
	"time"
	"w2learn/internal/dto"
	"w2learn/internal/model"
	"w2learn/internal/repository"
	"w2learn/internal/utils"
	"w2learn/pkg/logger"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ AuthService = (*authService)(nil)

type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) error
	Login(ctx context.Context, req *dto.LoginRequest) (string, error)
	Logout(ctx context.Context, tokenId string) error
}

type authService struct {
	userRepository repository.UserRepository
	redisClient    *redis.Client
}

func NewAuthService(userRepository repository.UserRepository, redisClient *redis.Client) AuthService {
	return &authService{
		userRepository: userRepository,
		redisClient:    redisClient,
	}
}

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) error {
	if req == nil {
		return errors.New("request is nil")
	}

	user, err := s.userRepository.GetByUsername(ctx, req.Username)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Failed to query user", zap.Error(err), zap.String("username", req.Username))
		return err
	}

	if user != nil {
		return errors.New("user already exists")
	}

	salt, err := utils.GenerateStringSalt(16)

	if err != nil || salt == "" {
		return errors.New("failed to generate salt")
	}

	user = &model.User{
		Username: req.Username,
		Password: utils.HashString(req.Password, salt),
		Salt:     salt,
		Status:   0,
		Habits:   nil,
	}

	err = s.userRepository.Create(ctx, user)

	if err != nil {
		return err
	}

	logger.Info("User created successfully",
		zap.String("username", user.Username),
		zap.Uint64("user_id", user.ID),
	)

	return nil
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest) (string, error) {
	if req == nil {
		return "", errors.New("request is nil")
	}

	user, err := s.userRepository.GetByUsername(ctx, req.Username)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Failed to query user", zap.Error(err), zap.String("username", req.Username))
		return "", err
	}

	if user == nil {
		return "", errors.New("user not found")
	}

	if !utils.VerifyString(req.Password, user.Salt, user.Password) {
		return "", errors.New("invalid password")
	}

	logger.Info("User logged in successfully",
		zap.String("username", user.Username),
		zap.Uint64("user_id", user.ID),
	)

	now := time.Now()

	token, err := utils.GenerateJwtToken(&dto.UserToken{
		UID:      user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 2)),
		},
	})

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Logout(ctx context.Context, tokenId string) error {
	return s.redisClient.Set(ctx, tokenId, true, 0).Err()
}
