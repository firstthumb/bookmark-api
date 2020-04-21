package bookmark

import (
	"bookmark-api/internal/bookmark/mocks"
	"bookmark-api/internal/entity"
	"bookmark-api/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func getFakeBookmark() entity.Bookmark {
	return entity.Bookmark{
		ID:   string(faker.RandomInt(0, 10000)),
		Name: faker.Lorem().Word(),
		Url:  faker.Internet().Url(),
		Tags: []string{faker.Hacker().Adjective()},
	}
}

func TestCreateRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	zapLogger := logger.NewLogger()

	mockRepository := mocks.NewMockRepository(ctrl)
	bookmarkService := NewService(mockRepository, zapLogger)

	api := NewApi(bookmarkService, zapLogger)

	r := gin.Default()

	api.RegisterHandlers(r.Group("/api"))

	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("CreateBookmarkSuccessfully", func(t *testing.T) {
		bookmark := getFakeBookmark()
		mockRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(bookmark, nil).Times(1)

		requestBody, _ := json.Marshal(CreateBookmarkRequest{
			Name: bookmark.Name,
			Url:  bookmark.Url,
			Tags: bookmark.Tags,
		})
		resp, err := http.Post(fmt.Sprintf("%s/api/bookmarks", ts.URL), "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		assert.Equal(t, 201, resp.StatusCode)

		var result BookmarkResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Expected bookmark response, got %v", err)
		}

		assert.Equal(t, bookmark.ID, result.ID)
		assert.Equal(t, bookmark.Name, result.Name)
		assert.Equal(t, bookmark.Url, result.Url)
		assert.Equal(t, bookmark.Tags, result.Tags)
	})

	t.Run("CreateBookmarkWithoutTagsSuccessfully", func(t *testing.T) {
		bookmark := getFakeBookmark()
		bookmark.Tags = nil
		mockRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(bookmark, nil).Times(1)

		requestBody, _ := json.Marshal(CreateBookmarkRequest{
			Name: bookmark.Name,
			Url:  bookmark.Url,
		})
		resp, err := http.Post(fmt.Sprintf("%s/api/bookmarks", ts.URL), "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		assert.Equal(t, 201, resp.StatusCode)

		var result BookmarkResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Expected bookmark response, got %v", err)
		}

		assert.Equal(t, bookmark.ID, result.ID)
		assert.Equal(t, bookmark.Name, result.Name)
		assert.Equal(t, bookmark.Url, result.Url)
		assert.Nil(t, bookmark.Tags)
	})
}

func TestUpdateRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	zapLogger := logger.NewLogger()

	mockRepository := mocks.NewMockRepository(ctrl)
	bookmarkService := NewService(mockRepository, zapLogger)

	api := NewApi(bookmarkService, zapLogger)

	r := gin.Default()

	api.RegisterHandlers(r.Group("/api"))

	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("UpdateBookmarkSuccessfully", func(t *testing.T) {
		bookmark := getFakeBookmark()
		mockRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(bookmark, nil).Times(1)

		requestBody, _ := json.Marshal(CreateBookmarkRequest{
			Name: bookmark.Name,
			Url:  bookmark.Url,
			Tags: bookmark.Tags,
		})

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/bookmarks/%s", ts.URL, bookmark.ID), strings.NewReader(string(requestBody)))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		assert.Equal(t, 200, resp.StatusCode)

		var result BookmarkResponse
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			t.Fatalf("Expected bookmark response, got %v", err)
		}

		assert.Equal(t, bookmark.ID, result.ID)
		assert.Equal(t, bookmark.Name, result.Name)
		assert.Equal(t, bookmark.Url, result.Url)
		assert.Equal(t, bookmark.Tags, result.Tags)
	})

}

func TestTagRoute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	zapLogger := logger.NewLogger()

	mockRepository := mocks.NewMockRepository(ctrl)
	bookmarkService := NewService(mockRepository, zapLogger)

	api := NewApi(bookmarkService, zapLogger)

	r := gin.Default()

	api.RegisterHandlers(r.Group("/api"))

	ts := httptest.NewServer(r)
	defer ts.Close()

	t.Run("AddTag", func(t *testing.T) {
		bookmark := getFakeBookmark()
		mockRepository.EXPECT().AddTag(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		tag := "test_tag"
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/bookmarks/%s/tags/%s", ts.URL, bookmark.ID, tag), strings.NewReader(""))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		assert.Equal(t, 201, resp.StatusCode)
	})

	t.Run("RemoveTag", func(t *testing.T) {
		bookmark := getFakeBookmark()
		mockRepository.EXPECT().RemoveTag(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		tag := "test_tag"
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/bookmarks/%s/tags/%s", ts.URL, bookmark.ID, tag), strings.NewReader(""))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		assert.Equal(t, 200, resp.StatusCode)
	})
}
