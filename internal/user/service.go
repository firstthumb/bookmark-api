package user

import (
	"bookmark-api/internal/entity"
	"context"
	"time"

	"go.uber.org/zap"
)

type Service interface {
	UpdateLastLogin(ctx context.Context, username string) (User, error)
}

type User struct {
	Username    string
	Method      string
	LastLoginAt time.Time
}

func (u *User) getEntity() entity.User {
	return entity.User{
		Username:    u.Username,
		Method:      u.Method,
		LastLoginAt: u.LastLoginAt,
	}
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

func (s *service) UpdateLastLogin(ctx context.Context, username string) (User, error) {
	logger := s.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	user, err := s.repo.Get(ctx, username)
	if err != nil {
		logger.Errorw("Could not find user", zap.String("Username", username))
		return User{}, err
	}

	user.LastLoginAt = time.Now()
	updatedUser, err := s.repo.Update(ctx, user)
	if err != nil {
		logger.Errorw("Failed to update last login")
		return User{}, err
	}

	return newUser(updatedUser), nil
}
