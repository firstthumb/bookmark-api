// repository.go
//go:generate mockgen -destination=mocks/repository_mock.go -package=mocks . Repository
package bookmark

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
	Create(ctx context.Context, bookmark entity.Bookmark) (entity.Bookmark, error)
	Get(ctx context.Context, id string) (entity.Bookmark, error)
	Update(ctx context.Context, bookmark entity.Bookmark) (entity.Bookmark, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db     *dynamo.DB
	logger *zap.Logger
}

func NewRepository(logger *zap.Logger) Repository {
	return &repository{db: db.GetDynamoDb(), logger: logger}
}

func (r *repository) Create(ctx context.Context, bookmark entity.Bookmark) (entity.Bookmark, error) {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	table := r.db.Table(db.TableBookmark)

	bookmark.ID = db.GenerateID()
	bookmark.CreatedAt = time.Now()
	bookmark.UpdatedAt = time.Now()
	err := table.Put(bookmark).Run()
	if err != nil {
		logger.Errorw("Failed to create bookmark")
		return entity.Bookmark{}, err
	}

	return bookmark, nil
}

func (r *repository) Get(ctx context.Context, id string) (entity.Bookmark, error) {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	table := r.db.Table(db.TableBookmark)

	var result entity.Bookmark
	err := table.Get("id", id).One(&result)
	if err != nil {
		logger.Errorw("Failed to get bookmark")
		switch err {
		case dynamo.ErrNotFound:
			return entity.Bookmark{}, errors.ErrNotFound
		default:
			return entity.Bookmark{}, err
		}
	}

	return result, nil
}

func (r *repository) Update(ctx context.Context, bookmark entity.Bookmark) (entity.Bookmark, error) {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	table := r.db.Table((db.TableBookmark))

	bookmark.UpdatedAt = time.Now()
	err := table.Put(bookmark).Run()
	if err != nil {
		logger.Errorw("Failed to update bookmark", zap.String("ID", bookmark.ID))
		return entity.Bookmark{}, err
	}

	return entity.Bookmark{}, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	table := r.db.Table((db.TableBookmark))

	err := table.Delete("id", id).Run()
	if err != nil {
		logger.Errorw("Failed to delete bookmark", zap.String("ID", id))
		return err
	}

	return nil
}
