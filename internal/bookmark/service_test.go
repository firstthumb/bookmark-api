package bookmark

import (
	"bookmark-api/internal/bookmark/mocks"
	"bookmark-api/internal/entity"
	"bookmark-api/pkg/logger"
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateBookmark(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockRepository(ctrl)
	s := NewService(mockRepository, logger.NewLogger())

	bookmark := entity.Bookmark{
		ID:        "ID",
		Name:      "Name",
		Url:       "Url",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	ctx := context.Background()

	mockRepository.EXPECT().Create(ctx, gomock.Any()).Return(bookmark, nil).Times(1)

	createdBookmark, err := s.Create(ctx, newBookmark(bookmark))
	assert.Nil(t, err)
	assert.NotNil(t, createdBookmark)
	assert.NotEmpty(t, createdBookmark.ID)
}
