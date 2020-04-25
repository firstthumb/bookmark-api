//+build wireinject

package di

import (
	"bookmark-api/internal/auth"
	"bookmark-api/internal/bookmark"
	"bookmark-api/pkg/logger"

	"github.com/google/wire"
)

var inject = wire.NewSet(logger.Inject, bookmark.Inject, auth.Inject)
var injectAuthorizer = wire.NewSet(logger.Inject, auth.Inject)

func CreateBookmarkApi() (bookmark.Api, error) {
	panic(wire.Build(inject))
}

func CreateAuthApi() (auth.Api, error) {
	panic(wire.Build(inject))
}

func CreateBookmarkService() (bookmark.Service, error) {
	panic(wire.Build(inject))
}

// func CreateAuthService() (auth.AuthService, error) {
// 	panic(wire.Build(injectAuthorizer))
// }

func CreateAuth() (*auth.Auth, error) {
	panic(wire.Build(inject))
}
