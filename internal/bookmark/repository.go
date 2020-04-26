// repository.go
//go:generate mockgen -destination=mocks/repository_mock.go -package=mocks . Repository
package bookmark

import (
	"context"
	"time"

	"github.com/guregu/dynamo"
	"github.com/thoas/go-funk"
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
	AddTag(ctx context.Context, tag entity.Tag) error
	RemoveTag(ctx context.Context, tag entity.Tag) error
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

	tx := r.db.WriteTx()

	// Create Bookmark
	tableBookmark := r.db.Table(db.GetTableBookmark())

	bookmark.ID = db.GenerateID()
	bookmark.CreatedAt = time.Now()
	bookmark.UpdatedAt = time.Now()
	tx.Put(tableBookmark.Put(bookmark))

	// Create Tags
	tableTag := r.db.Table(db.GetTableTag())
	for _, tag := range bookmark.Tags {
		tx.Put(tableTag.Put(entity.Tag{Tag: tag, BookmarkID: bookmark.ID}))
	}

	err := tx.Run()
	if err != nil {
		logger.Errorw("Failed to create bookmark", zap.Error(err))
		return entity.Bookmark{}, err
	}

	return bookmark, nil
}

func (r *repository) Get(ctx context.Context, id string) (entity.Bookmark, error) {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	table := r.db.Table(db.GetTableBookmark())

	var result entity.Bookmark
	err := table.Get("id", id).One(&result)
	if err != nil {
		logger.Errorw("Failed to get bookmark", zap.Error(err))
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

	freshBookmark, err := r.Get(ctx, bookmark.ID)
	if err != nil {
		return entity.Bookmark{}, err
	}

	table := r.db.Table((db.GetTableBookmark()))

	// Don't update tags here
	freshBookmark.Name = bookmark.Name
	freshBookmark.Url = bookmark.Url
	freshBookmark.UpdatedAt = time.Now()
	err = table.Put(freshBookmark).Run()
	if err != nil {
		logger.Errorw("Failed to update bookmark", zap.String("ID", freshBookmark.ID), zap.Error(err))
		return entity.Bookmark{}, err
	}

	return entity.Bookmark{}, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	bookmark, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	tx := r.db.WriteTx()

	tableBookmark := r.db.Table((db.GetTableBookmark()))
	tableTag := r.db.Table((db.GetTableTag()))

	tx.Delete(tableBookmark.Delete("id", id))
	for _, tag := range bookmark.Tags {
		tx.Delete(tableTag.Delete("tag", tag).Range("bookmark_id", id))
	}

	err = tx.Run()
	if err != nil {
		logger.Errorw("Failed to delete bookmark", zap.String("ID", id), zap.Error(err))
		return err
	}

	return nil
}

func (r *repository) AddTag(ctx context.Context, tag entity.Tag) error {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	bookmark, err := r.Get(ctx, tag.BookmarkID)
	if err != nil {
		return err
	}

	if funk.Contains(bookmark.Tags, tag.Tag) {
		logger.Errorw("Already has tag", zap.String("BookmarkID", tag.BookmarkID), zap.String("Tag", tag.Tag), zap.Error(err))
		return errors.ErrAlreadyExist
	}

	tx := r.db.WriteTx()

	// Update Bookmark
	tableBookmark := r.db.Table(db.GetTableBookmark())
	bookmark.Tags = append(bookmark.Tags, tag.Tag)
	tx.Put(tableBookmark.Put(bookmark))

	// Add Tag
	tableTag := r.db.Table(db.GetTableTag())
	tx.Put(tableTag.Put(tag))

	err = tx.Run()
	if err != nil {
		logger.Errorw("Failed to create tag", zap.Error(err))
		return err
	}

	return nil
}

func (r *repository) RemoveTag(ctx context.Context, tag entity.Tag) error {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	bookmark, err := r.Get(ctx, tag.BookmarkID)
	if err != nil {
		return err
	}

	if !funk.Contains(bookmark.Tags, tag.Tag) {
		logger.Errorw("Bookmark has not tag", zap.String("BookmarkID", tag.BookmarkID), zap.String("Tag", tag.Tag), zap.Error(err))
		return errors.ErrInvalidParam
	}

	tx := r.db.WriteTx()

	// Update Bookmark
	tableBookmark := r.db.Table(db.GetTableBookmark())
	bookmark.Tags = funk.FilterString(bookmark.Tags, func(s string) bool { return s != tag.Tag })
	tx.Put(tableBookmark.Put(bookmark))

	// Add Tag
	tableTag := r.db.Table(db.GetTableTag())
	tx.Put(tableTag.Put(tag))

	err = tx.Run()
	if err != nil {
		logger.Errorw("Failed to create tag", zap.Error(err))
		return err
	}

	return nil
}
