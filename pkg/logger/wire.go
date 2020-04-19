package logger

import (
	"github.com/google/wire"
)

var Inject = wire.NewSet(NewLogger)
