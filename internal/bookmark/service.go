package bookmark

import (
	"bookmark-api/internal/entity"
	"context"

	"github.com/thoas/go-funk"
	"go.uber.org/zap"
)

type Service interface {
	Create(ctx context.Context, bookmark Bookmark) (Bookmark, error)
	Get(ctx context.Context, username, bookmarkId string) (Bookmark, error)
	Update(ctx context.Context, bookmark Bookmark) (Bookmark, error)
	Delete(ctx context.Context, username, bookmarkId string) error
	SearchByName(ctx context.Context, username, name string) ([]Bookmark, error)
	AddTag(ctx context.Context, username, bookmarkId, tag string) error
	RemoveTag(ctx context.Context, username, bookmarkId, tag string) error
}

type Bookmark struct {
	ID       string
	Username string
	Name     string
	Url      string
	Tags     []string
}

func (b *Bookmark) getEntity() entity.Bookmark {
	return entity.Bookmark{
		ID:       b.ID,
		Username: b.Username,
		Name:     b.Name,
		Url:      b.Url,
		Tags:     b.Tags,
	}
}

func newBookmark(bookmark entity.Bookmark) Bookmark {
	return Bookmark{
		ID:       bookmark.GetBookmarkId(),
		Username: bookmark.GetUsername(),
		Name:     bookmark.Name,
		Url:      bookmark.Url,
		Tags:     bookmark.Tags,
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

func (s *service) Get(ctx context.Context, username, bookmarkId string) (Bookmark, error) {
	logger := s.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	result, err := s.repo.Get(ctx, username, bookmarkId)
	if err != nil {
		logger.Errorw("Failed to fetch", zap.String("ID", bookmarkId))
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

func (s *service) Delete(ctx context.Context, username, bookmarkId string) error {
	logger := s.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	err := s.repo.Delete(ctx, username, bookmarkId)
	if err != nil {
		logger.Errorw("Failed to delete", zap.String("ID", bookmarkId))
		return err
	}

	return nil
}

func (s *service) SearchByName(ctx context.Context, username, name string) ([]Bookmark, error) {
	logger := s.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	result, err := s.repo.SearchByName(ctx, username, name)
	if err != nil {
		logger.Errorw("Failed to search bookmark by name", zap.Error(err))
		return []Bookmark{}, err
	}

	return funk.Map(result, func(b entity.Bookmark) Bookmark {
		return newBookmark(b)
	}).([]Bookmark), nil
}

func (s *service) AddTag(ctx context.Context, username, bookmarkId, tag string) error {
	logger := s.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	err := s.repo.AddTag(ctx, username, bookmarkId, tag)
	if err != nil {
		logger.Errorw("Failed to add tag", zap.String("BookmarkID", bookmarkId), zap.String("Tag", tag))
		return err
	}

	return nil
}

func (s *service) RemoveTag(ctx context.Context, username, bookmarkId, tag string) error {
	logger := s.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	err := s.repo.RemoveTag(ctx, username, bookmarkId, tag)
	if err != nil {
		logger.Errorw("Failed to remove tag", zap.String("BookmarkID", bookmarkId), zap.String("Tag", tag))
		return err
	}

	return nil
}
