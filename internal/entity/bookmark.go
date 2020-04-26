package entity

import (
	"time"
)

type Bookmark struct {
	ID        string    `json:"id" dynamo:"id"`
	Name      string    `json:"name" dynamo:"name"`
	Url       string    `json:"url" dynamo:"url"`
	Tags      []string  `json:"tags" dynamo:"tags"`
	CreatedAt time.Time `json:"created_at" dynamo:"created_at"`
	UpdatedAt time.Time `json:"updated_at" dynamo:"updated_at"`
}
