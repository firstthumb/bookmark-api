package bookmark

import (
	"github.com/google/wire"
)

var Inject = wire.NewSet(NewApi, NewRepository, NewService)
