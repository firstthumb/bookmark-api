package auth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AuthMethod string

const (
	Google AuthMethod = "google"
)

type AuthUser struct {
	Username string
	Method   AuthMethod
}

type Claim struct {
	Username string `json:"username"`
	Method   string `json:"method"`
}

type AuthService interface {
	SigninHandler(ctx *gin.Context)
	GetSigninURL(state string) string
	VerifyToken(token *oauth2.Token) (*AuthUser, error)
	AuthCallbackMiddleware() gin.HandlerFunc
}
