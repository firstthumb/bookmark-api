package auth

import (
	"errors"
	"log"
	"os"
	"time"

	verifier "github.com/gbrlsnchs/jwt/v3"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	jwt "github.com/appleboy/gin-jwt/v2"

	"bookmark-api/internal/user"
)

func NewAuth(google *GoogleOAuth, userService user.Service, logger *zap.Logger) *Auth {
	return &Auth{Google: google, userService: userService, logger: logger}
}

type Auth struct {
	Google      *GoogleOAuth
	userService user.Service
	logger      *zap.Logger
}

func (a *Auth) Session() gin.HandlerFunc {
	store := sessions.NewCookieStore([]byte("sessionKey"))

	return sessions.Sessions("auth_session", store)
}

func (a *Auth) VerifyToken(token string) (Claim, error) {
	logger := a.logger.Sugar()
	defer func() {
		_ = logger.Sync()
	}()

	hs := verifier.NewHS256([]byte(os.Getenv("OAUTH_KEY")))
	var payload Claim
	_, err := verifier.Verify([]byte(token), hs, &payload)
	if err != nil {
		logger.Errorw("Failed to verify", zap.String("Token", token), zap.Error(err))
		return Claim{}, err
	}

	logger.Infow("Verify token", zap.String("Token", token), zap.String("Username", payload.Username), zap.String("Method", payload.Method))

	return payload, nil
}

func (a *Auth) AuthMiddleware() *jwt.GinJWTMiddleware {
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       os.Getenv("OAUTH_REALM"),
		Key:         []byte(os.Getenv("OAUTH_KEY")),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour * 24 * 7,
		IdentityKey: "username",

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			logger := a.logger.Sugar()
			defer func() {
				_ = logger.Sync()
			}()

			if v, ok := data.(*AuthUser); ok {
				return jwt.MapClaims{
					"username": v.Username,
					"method":   v.Method,
				}
			}

			logger.Errorw("Wrong token claim")

			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			if claims["username"] == nil || claims["method"] == nil {
				return nil
			}

			username := claims["username"].(string)
			method := claims["method"].(string)

			authUser := &AuthUser{
				Username: username,
				Method:   AuthMethod(method),
			}

			// Save authUser in the context
			c.Set("user", authUser)

			return authUser
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			token, exist := c.Get("token")

			if !exist {
				return nil, errors.New("Token does not exist")
			}

			authUser, err := a.Google.VerifyToken(token.(*oauth2.Token))
			if err != nil {
				return nil, err
			}

			// Update last login date
			_, _ = a.userService.UpdateLastLogin(c.Request.Context(), authUser.Username, string(authUser.Method))

			return authUser, nil
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			if u, ok := data.(*AuthUser); ok && u.Method != "" {
				return a.Google.Authorize(c, u)
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    "error",
				"message": "failed",
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return middleware
}
