package auth

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewApi(auth *Auth, logger *zap.Logger) Api {
	return &resource{auth, logger}
}

type resource struct {
	auth   *Auth
	logger *zap.Logger
}

type Api interface {
	Session() gin.HandlerFunc
	RegisterSigninHandlers(rg *gin.RouterGroup)
	RegisterAuthHandlers(rg *gin.RouterGroup)
	GetAuthMiddleware() *jwt.GinJWTMiddleware
}

func (r *resource) Session() gin.HandlerFunc {
	return r.auth.Session()
}

func (r *resource) GetAuthMiddleware() *jwt.GinJWTMiddleware {
	return r.auth.AuthMiddleware()
}

func (r *resource) RegisterSigninHandlers(rg *gin.RouterGroup) {
	rg.GET("/signin/google", r.auth.Google.SigninHandler)
	rg.GET("/callback", r.auth.Google.AuthCallbackMiddleware(), r.GetAuthMiddleware().LoginHandler)
}

func (r *resource) RegisterAuthHandlers(rg *gin.RouterGroup) {
	auth := rg.Group("/auth", r.GetAuthMiddleware().MiddlewareFunc())
	{
		auth.POST("/check", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"Status": "Success", "Message": "Token is valid"})
		})
		auth.POST("/refresh", r.GetAuthMiddleware().RefreshHandler)
	}
}
