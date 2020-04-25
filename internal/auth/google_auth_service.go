package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	googleAuthIDTokenVerifier "github.com/futurenda/google-auth-id-token-verifier"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const stateKey = "state"

type GoogleUser struct {
	Sub           string `json:"sub"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func randToken() string {
	buffer := make([]byte, 32)
	_, _ = rand.Read(buffer)
	return base64.StdEncoding.EncodeToString(buffer)
}

func NewGoogleOAuth() *GoogleOAuth {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	return &GoogleOAuth{
		conf: conf,
	}
}

type GoogleOAuth struct {
	conf *oauth2.Config
}

func (o *GoogleOAuth) SigninHandler(ctx *gin.Context) {
	state := randToken()
	session := sessions.Default(ctx)
	session.Set(stateKey, state)
	_ = session.Save()

	ctx.Redirect(http.StatusFound, o.GetSigninURL(state))
}

func (o *GoogleOAuth) GetSigninURL(state string) string {
	return o.conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (o *GoogleOAuth) VerifyToken(token *oauth2.Token) (*AuthUser, error) {
	idToken := token.Extra("id_token").(string)
	verifier := googleAuthIDTokenVerifier.Verifier{}
	err := verifier.VerifyIDToken(idToken, []string{
		os.Getenv("GOOGLE_CLIENT_ID"),
	})
	if err != nil {
		return nil, err
	}

	claimSet, err := googleAuthIDTokenVerifier.Decode(idToken)
	if err != nil {
		return nil, err
	}

	return &AuthUser{Username: claimSet.Email, Method: Google}, nil
}

func (o *GoogleOAuth) Authorize(c *gin.Context, user *AuthUser) bool {
	return true
}

func (o *GoogleOAuth) AuthCallbackMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		retrievedState := session.Get(stateKey)
		session.Delete(stateKey)

		if retrievedState != ctx.Query(stateKey) {
			_ = ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
			return
		}

		tok, err := o.conf.Exchange(oauth2.NoContext, ctx.Query("code"))
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		client := o.conf.Client(oauth2.NoContext, tok)
		email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		defer email.Body.Close()
		data, err := ioutil.ReadAll(email.Body)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var user GoogleUser
		err = json.Unmarshal(data, &user)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.Set("user", user)
		ctx.Set("token", tok)
	}
}
