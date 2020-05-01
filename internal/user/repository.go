// repository.go
//go:generate mockgen -destination=mocks/repository_mock.go -package=mocks . Repository
package user

import (
	"context"
	"time"

	"github.com/guregu/dynamo"
	"go.uber.org/zap"

	"bookmark-api/internal/entity"
	"bookmark-api/internal/errors"
	"bookmark-api/pkg/db"
)

type Repository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Get(ctx context.Context, username string) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
}

type repository struct {
	db     *dynamo.DB
	logger *zap.Logger
}

func NewRepository(logger *zap.Logger) Repository {
	return &repository{db: db.GetDynamoDb(), logger: logger}
}

func (r *repository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	tx := r.db.WriteTx()

	tableUser := r.db.Table(db.GetTableUser())

	user.LastLoginAt = time.Now()
	tx.Put(tableUser.Put(user))

	err := tx.Run()
	if err != nil {
		logger.Errorw("Failed to create user", zap.Error(err))
		return entity.User{}, err
	}

	return user, nil
}

func (r *repository) Get(ctx context.Context, username string) (entity.User, error) {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	table := r.db.Table(db.GetTableUser())

	// TODO: Support multiple methods
	var result entity.User
	err := table.Get("username", username).One(&result)
	if err != nil {
		logger.Errorw("Failed to get user", zap.Error(err))
		switch err {
		case dynamo.ErrNotFound:
			return entity.User{}, errors.ErrNotFound
		default:
			return entity.User{}, err
		}
	}

	return result, nil
}

func (r *repository) Update(ctx context.Context, user entity.User) (entity.User, error) {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	freshUser, err := r.Get(ctx, user.Username)
	if err != nil {
		return entity.User{}, err
	}

	table := r.db.Table((db.GetTableUser()))
	freshUser.LastLoginAt = user.LastLoginAt
	err = table.Put(freshUser).Run()
	if err != nil {
		logger.Errorw("Failed to update user", zap.String("Username", freshUser.Username), zap.Error(err))
		return entity.User{}, err
	}

	return entity.User{}, nil
}
