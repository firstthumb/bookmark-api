package user

import (
	"github.com/google/wire"
)

var Inject = wire.NewSet(NewRepository, NewService)
