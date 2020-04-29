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
	Get(ctx context.Context, username, id string) (entity.Bookmark, error)
	Update(ctx context.Context, bookmark entity.Bookmark) (entity.Bookmark, error)
	Delete(ctx context.Context, username, id string) error
	SearchByName(ctx context.Context, username, name string) ([]entity.Bookmark, error)
	AddTag(ctx context.Context, username, bookmarkId, tag string) error
	RemoveTag(ctx context.Context, username, bookmarkId, tag string) error
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

	tableBookmark := r.db.Table(db.GetTableBookmark())

	// Create Bookmark
	bookmark.ID = db.GenerateID()
	bookmark.CreatedAt = time.Now()
	bookmark.UpdatedAt = time.Now()
	tx.Put(tableBookmark.Put(bookmark.GetEntity()))

	// Create SearchByName
	tx.Put(tableBookmark.Put(bookmark.GetSearchByName()))

	// Create SearchByTag
	for _, searchByTag := range bookmark.GetSearchByTag() {
		tx.Put(tableBookmark.Put(searchByTag))
	}

	err := tx.Run()
	if err != nil {
		logger.Errorw("Failed to create bookmark", zap.Error(err))
		return entity.Bookmark{}, err
	}

	return bookmark, nil
}

func (r *repository) Get(ctx context.Context, username, bookmarkId string) (entity.Bookmark, error) {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	table := r.db.Table(db.GetTableBookmark())

	hashId, rangeId := entity.GetSearchKeyByID(username, bookmarkId)
	var result entity.Bookmark
	err := table.Get("id", hashId).
		Range("range", "EQ", rangeId).
		One(&result)
	if err != nil {
		logger.Errorw("Failed to get bookmark", zap.String("Username", username), zap.String("ID", bookmarkId), zap.Error(err))
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

	tx := r.db.WriteTx()

	freshBookmark, err := r.Get(ctx, bookmark.Username, bookmark.ID)
	if err != nil {
		return entity.Bookmark{}, err
	}

	table := r.db.Table((db.GetTableBookmark()))

	// Update URL only
	freshBookmark.Url = bookmark.Url
	freshBookmark.UpdatedAt = time.Now()
	tx.Put(table.Put(freshBookmark))

	err = tx.Run()
	if err != nil {
		logger.Errorw("Failed to update bookmark", zap.String("ID", freshBookmark.ID), zap.Error(err))
		return entity.Bookmark{}, err
	}

	return entity.Bookmark{}, nil
}

func (r *repository) Delete(ctx context.Context, username, bookmarkId string) error {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	// TODO: Implement this

	return nil
}

func (r *repository) SearchByName(ctx context.Context, username string, name string) ([]entity.Bookmark, error) {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	tableBookmark := r.db.Table(db.GetTableBookmark())

	var result []entity.Bookmark

	// Search by name
	hashId, rangeId := entity.GetSearchKeyByName(username, name)
	var searchByNameResult []entity.BookmarkSearchByName
	err := tableBookmark.Get("id", hashId).
		Range("range", "BEGINS_WITH", rangeId).
		All(&searchByNameResult)

	if err != nil {
		logger.Errorw("Failed to search bookmark", zap.String("HashId", hashId), zap.String("RangeId", rangeId), zap.Error(err))
		return []entity.Bookmark{}, err
	}

	// Fetch every bookmark by ID
	for _, bookmarkName := range searchByNameResult {
		hashId, rangeId = entity.GetSearchKeyByID(username, bookmarkName.GetBookmarkId())
		var bookmark entity.Bookmark
		err = tableBookmark.Get("id", hashId).
			Range("range", "EQ", rangeId).One(&bookmark)

		if err != nil {
			logger.Errorw("Could not fetch bookmark", zap.String("HashId", hashId), zap.String("RangeId", rangeId))
		} else {
			result = append(result, bookmark)
		}
	}

	return result, nil
}

func (r *repository) AddTag(ctx context.Context, username, bookmarkId, tag string) error {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	bookmark, err := r.Get(ctx, username, bookmarkId)
	if err != nil {
		return err
	}

	if funk.Contains(bookmark.Tags, tag) {
		logger.Errorw("Already has tag", zap.String("BookmarkID", bookmarkId), zap.String("Tag", tag), zap.Error(err))
		return errors.ErrAlreadyExist
	}

	tx := r.db.WriteTx()

	// Update Bookmark
	table := r.db.Table(db.GetTableBookmark())
	bookmark.Tags = append(bookmark.Tags, tag)
	tx.Put(table.Put(bookmark))

	// Add Tag
	tx.Put(table.Put(entity.NewBookmarkSearchByTag(username, bookmarkId, tag)))

	err = tx.Run()
	if err != nil {
		logger.Errorw("Failed to create tag", zap.Error(err))
		return err
	}

	return nil
}

func (r *repository) RemoveTag(ctx context.Context, username, bookmarkId, tag string) error {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	bookmark, err := r.Get(ctx, username, bookmarkId)
	if err != nil {
		return err
	}

	if !funk.Contains(bookmark.Tags, tag) {
		logger.Errorw("Bookmark has not tag", zap.String("BookmarkID", bookmarkId), zap.String("Tag", tag), zap.Error(err))
		return errors.ErrInvalidParam
	}

	tx := r.db.WriteTx()

	// Update Bookmark
	table := r.db.Table(db.GetTableBookmark())
	bookmark.Tags = funk.FilterString(bookmark.Tags, func(s string) bool { return s != tag })
	tx.Put(table.Put(bookmark))

	// Delete Tag
	searchTag := entity.NewBookmarkSearchByTag(username, bookmarkId, tag)
	tx.Delete(table.Delete("id", searchTag.Username).Range("range", searchTag.Tag))

	err = tx.Run()
	if err != nil {
		logger.Errorw("Failed to create tag", zap.Error(err))
		return err
	}

	return nil
}
