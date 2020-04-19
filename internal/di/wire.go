//+build wireinject

package di

import (
	"bookmark-api/internal/bookmark"
	"bookmark-api/pkg/logger"

	"github.com/google/wire"
)

var inject = wire.NewSet(logger.Inject, bookmark.Inject)

func CreateBookmarkApi() (bookmark.Api, error) {
	panic(wire.Build(inject))
}

func CreateBookmarkService() (bookmark.Service, error) {
	panic(wire.Build(inject))
}
