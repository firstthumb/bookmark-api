//+build wireinject

package di

import (
	"bookmark-api/internal/auth"
	"bookmark-api/internal/bookmark"
	"bookmark-api/internal/user"
	"bookmark-api/pkg/logger"

	"github.com/google/wire"
)

var inject = wire.NewSet(logger.Inject, bookmark.Inject, auth.Inject, user.Inject)
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

func CreateAuth() (*auth.Auth, error) {
	panic(wire.Build(inject))
}
