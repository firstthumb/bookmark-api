package bookmark

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"go.uber.org/zap"

	"bookmark-api/internal/errors"
	"bookmark-api/pkg/utils"
)

func NewApi(service Service, logger *zap.Logger) Api {
	return &resource{service, logger}
}

type Api interface {
	RegisterHandlers(rg *gin.RouterGroup)
}

func (r *resource) RegisterHandlers(rg *gin.RouterGroup) {
	rg.POST("/bookmarks", r.create)
	rg.GET("/bookmarks/:id", r.get)
	rg.PUT("/bookmarks/:id", r.update)
	rg.DELETE("/bookmarks/:id", r.delete)
	rg.POST("/bookmarks/:id/tags/:tag", r.addTag)
	rg.DELETE("/bookmarks/:id/tags/:tag", r.removeTag)

	rg.GET("/bookmarks", r.searchByName)
}

type CreateBookmarkRequest struct {
	Name string   `json:"name"`
	Url  string   `json:"url"`
	Tags []string `json:"tags"`
}

type UpdateBookmarkRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type BookmarkResponse struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Url  string   `json:"url"`
	Tags []string `json:"tags"`
}

type resource struct {
	service Service
	logger  *zap.Logger
}

func (r *resource) create(c *gin.Context) {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	authUser := utils.GetCurrentUser(c)

	request := CreateBookmarkRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Could not bind payload")
		c.JSON(http.StatusBadRequest, errors.BadRequest("Payload is in wrong format"))
		return
	}

	result, err := r.service.Create(c.Request.Context(), Bookmark{
		Name:     request.Name,
		Username: authUser.Username,
		Url:      request.Url,
		Tags:     request.Tags,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServerError("Failed to create bookmark"))
		return
	}

	response := &BookmarkResponse{
		ID:   result.ID,
		Name: result.Name,
		Url:  result.Url,
		Tags: result.Tags,
	}
	c.JSON(http.StatusCreated, response)
}

func (r *resource) get(c *gin.Context) {
	id := c.Param("id")
	if len(id) < 0 {
		c.JSON(http.StatusBadRequest, errors.BadRequest("Parameter(id) is missing"))
		return
	}

	authUser := utils.GetCurrentUser(c)
	result, err := r.service.Get(c.Request.Context(), authUser.Username, id)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, errors.NotFound("Not found"))
		default:
			c.JSON(http.StatusInternalServerError, errors.InternalServerError("Failed to get bookmark"))
		}
		return
	}

	response := &BookmarkResponse{
		ID:   result.ID,
		Name: result.Name,
		Url:  result.Url,
		Tags: result.Tags,
	}
	c.JSON(http.StatusOK, response)
}

func (r *resource) update(c *gin.Context) {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	id := c.Param("id")
	if len(id) < 0 {
		c.JSON(http.StatusBadRequest, errors.BadRequest("Parameter(id) is missing"))
		return
	}

	request := UpdateBookmarkRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Could not bind payload")
		c.JSON(http.StatusBadRequest, errors.BadRequest("Payload is in wrong format"))
		return
	}

	authUser := utils.GetCurrentUser(c)
	result, err := r.service.Update(c.Request.Context(), Bookmark{
		Username: authUser.Username,
		ID:       id,
		Name:     request.Name,
		Url:      request.Url,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServerError("Failed to update bookmark"))
		return
	}

	response := &BookmarkResponse{
		ID:   result.ID,
		Name: result.Name,
		Url:  result.Url,
		Tags: result.Tags,
	}
	c.JSON(http.StatusOK, response)
}

func (r *resource) delete(c *gin.Context) {
	id := c.Param("id")
	if len(id) < 0 {
		c.JSON(http.StatusBadRequest, errors.BadRequest("Parameter(id) is missing"))
		return
	}

	authUser := utils.GetCurrentUser(c)
	err := r.service.Delete(c.Request.Context(), authUser.Username, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServerError("Failed to delete bookmark"))
		return
	}

	c.Status(http.StatusOK)
}

func (r *resource) searchByName(c *gin.Context) {
	query := c.Query("query")
	if len(query) < 0 {
		c.JSON(http.StatusBadRequest, errors.BadRequest("Query is missing"))
		return
	}

	authUser := utils.GetCurrentUser(c)
	result, err := r.service.SearchByName(c.Request.Context(), authUser.Username, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServerError("Failed to delete bookmark"))
		return
	}

	response := funk.Map(result, func(bookmark Bookmark) BookmarkResponse {
		return BookmarkResponse{
			ID:   bookmark.ID,
			Name: bookmark.Name,
			Url:  bookmark.Url,
			Tags: bookmark.Tags,
		}
	})

	c.JSON(http.StatusOK, response)
}

func (r *resource) addTag(c *gin.Context) {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	bookmarkId := c.Param("id")
	tag := c.Param("tag")

	authUser := utils.GetCurrentUser(c)
	err := r.service.AddTag(c.Request.Context(), authUser.Username, bookmarkId, tag)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServerError("Failed to add tag"))
		return
	}

	c.Status(http.StatusCreated)
}

func (r *resource) removeTag(c *gin.Context) {
	logger := r.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	bookmarkId := c.Param("id")
	tag := c.Param("tag")

	authUser := utils.GetCurrentUser(c)
	err := r.service.RemoveTag(c.Request.Context(), authUser.Username, bookmarkId, tag)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServerError("Failed to remove tag"))
		return
	}

	c.Status(http.StatusOK)
}
