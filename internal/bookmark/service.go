package bookmark

import (
	"bookmark-api/internal/entity"
	"context"

	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, bookmark Bookmark) (Bookmark, error)
	Get(ctx context.Context, id string) (Bookmark, error)
	Update(ctx context.Context, bookmark Bookmark) (Bookmark, error)
	Delete(ctx context.Context, id string) error
}

type Bookmark struct {
	ID   string
	Name string
	Url  string
}

func (b *Bookmark) getEntity() entity.Bookmark {
	return entity.Bookmark{
		ID:   b.ID,
		Name: b.Name,
		Url:  b.Url,
	}
}

func newBookmark(bookmark entity.Bookmark) Bookmark {
	return Bookmark{
		ID:   bookmark.ID,
		Name: bookmark.ID,
		Url:  bookmark.Url,
	}
}

type service struct {
	repo   Repository
	logger *zap.Logger
}

func NewService(repo Repository, logger *zap.Logger) Service {
	return &service{repo, logger}
}

func (s *service) Create(ctx context.Context, bookmark Bookmark) (Bookmark, error) {
	logger := s.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	entityBookmark := bookmark.getEntity()
	createdBookmark, err := s.repo.Create(ctx, entityBookmark)
	if err != nil {
		logger.Errorw("Failed to create")
		return Bookmark{}, err
	}

	return newBookmark(createdBookmark), nil
}

func (s *service) Get(ctx context.Context, id string) (Bookmark, error) {
	logger := s.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	result, err := s.repo.Get(ctx, id)
	if err != nil {
		logger.Errorw("Failed to fetch", zap.String("ID", id))
		return Bookmark{}, err
	}

	return newBookmark(result), nil
}

func (s *service) Update(ctx context.Context, bookmark Bookmark) (Bookmark, error) {
	logger := s.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	updatedBookmark, err := s.repo.Update(ctx, bookmark.getEntity())
	if err != nil {
		logger.Errorw("Failed to update", zap.String("ID", bookmark.ID))
		return Bookmark{}, err
	}

	return newBookmark(updatedBookmark), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	logger := s.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	err := s.repo.Delete(ctx, id)
	if err != nil {
		logger.Errorw("Failed to delete", zap.String("ID", id))
		return err
	}

	return nil
}
