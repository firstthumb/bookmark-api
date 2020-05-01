package session

import (
	"bookmark-api/internal/auth"

	"github.com/gin-gonic/gin"
)

func GetCurrentUser(ctx *gin.Context) *auth.AuthUser {
	authUser, exist := ctx.Get("user")
	if exist {
		return authUser.(*auth.AuthUser)
	}

	return nil
}
