package entity

import (
	"time"
)

type User struct {
	Username    string    `json:"username" dynamo:"username"`
	Method      string    `json:"method" dynamo:"method"`
	LastLoginAt time.Time `json:"last_login_at" dynamo:"last_login_at"`
}
