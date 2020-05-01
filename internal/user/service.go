package user

import (
	"bookmark-api/internal/entity"
	"context"
	"time"

	"go.uber.org/zap"
)

type Service interface {
	CreateUser(ctx context.Context, username, method string) (User, error)
	UpdateLastLogin(ctx context.Context, username, method string) (User, error)
}

type User struct {
	Username    string
	Method      string
	LastLoginAt time.Time
}

func newUser(user entity.User) User {
	return User{
		Username:    user.Username,
		Method:      user.Method,
		LastLoginAt: user.LastLoginAt,
	}
}

type service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {
	return &service{repo, logger}
}

func (s *service) CreateUser(ctx context.Context, username, method string) (User, error) {
	logger := s.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	user, err := s.repo.Create(ctx, entity.User{
		Username:    username,
		Method:      method,
		LastLoginAt: time.Now(),
	})

	if err != nil {
		logger.Errorw("Could not create user", zap.String("Username", username))
		return User{}, err
	}

	return newUser(user), nil
}

func (s *service) UpdateLastLogin(ctx context.Context, username, method string) (User, error) {
	logger := s.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	user, err := s.repo.Get(ctx, username)
	if err != nil {
		// Create new user
		return s.CreateUser(ctx, username, method)
	}

	user.LastLoginAt = time.Now()
	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		logger.Errorw("Failed to update last login")
		return User{}, err
	}

	return newUser(updatedUser), nil
}
