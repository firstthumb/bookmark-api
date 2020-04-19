package utils

import (
	"os"
	"strconv"
)

func IsOffline() bool {
	r, err := strconv.ParseBool(os.Getenv("IS_OFFLINE"))
	if err != nil {
		return false
	}

	return r
}
