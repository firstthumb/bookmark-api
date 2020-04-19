package bookmark

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"bookmark-api/internal/errors"
)

func NewApi(service Service, logger *zap.Logger) Api {
	return &resource{service, logger}
}

type Api interface {
	RegisterHandlers(rg *gin.RouterGroup)
}

func (r *resource) RegisterHandlers(rg *gin.RouterGroup) {
	rg.POST("/bookmarks/", r.create)
	rg.GET("/bookmarks/:id", r.get)
	rg.PUT("/bookmarks/:id", r.update)
	rg.DELETE("/bookmarks/:id", r.delete)
}

type CreateBookmarkRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type UpdateBookmarkRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
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

	request := CreateBookmarkRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Could not bind payload")
		c.JSON(http.StatusBadRequest, errors.BadRequest("Payload is in wrong format"))
		return
	}

	result, err := r.service.Create(c.Request.Context(), Bookmark{
		Name: request.Name,
		Url:  request.Url,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServerError("Failed to create bookmark"))
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (r *resource) get(c *gin.Context) {
	id := c.Param("id")
	if len(id) < 0 {
		c.JSON(http.StatusBadRequest, errors.BadRequest("Parameter(id) is missing"))
		return
	}

	result, err := r.service.Get(c.Request.Context(), id)
	if err != nil {
		switch err {
		case errors.ErrNotFound:
			c.JSON(http.StatusNotFound, errors.NotFound("Not found"))
		default:
			c.JSON(http.StatusInternalServerError, errors.InternalServerError("Failed to get bookmark"))
		}
		return
	}

	c.JSON(http.StatusOK, result)
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

	result, err := r.service.Update(c.Request.Context(), Bookmark{
		ID:   id,
		Name: request.Name,
		Url:  request.Url,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServerError("Failed to update bookmark"))
		return
	}

	c.JSON(http.StatusOK, result)
}

func (r *resource) delete(c *gin.Context) {
	id := c.Param("id")
	if len(id) < 0 {
		c.JSON(http.StatusBadRequest, errors.BadRequest("Parameter(id) is missing"))
		return
	}

	err := r.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalServerError("Failed to delete bookmark"))
		return
	}

	c.Status(http.StatusOK)
}
