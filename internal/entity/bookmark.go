package entity

import (
	"time"
)

type Bookmark struct {
	ID        string    `json:"id" dynamo:"id"`
	Name      string    `json:"name" dynamo:"name"`
	Url       string    `json:"url" dynamo:"-"`
	CreatedAt time.Time `json:"created_at" dynamo:"-"`
	UpdatedAt time.Time `json:"updated_at" dynamo:"-"`
}
