package auth

import (
	"github.com/google/wire"
)

var Inject = wire.NewSet(NewApi, NewAuth, NewGoogleOAuth)
