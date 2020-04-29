package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/thoas/go-funk"
)

type Bookmark struct {
	Username  string    `json:"username" dynamo:"id"`
	ID        string    `json:"id" dynamo:"range"`
	Name      string    `json:"name" dynamo:"name"`
	Url       string    `json:"url" dynamo:"url"`
	Tags      []string  `json:"tags" dynamo:"tags"`
	CreatedAt time.Time `json:"created_at" dynamo:"created_at"`
	UpdatedAt time.Time `json:"updated_at" dynamo:"updated_at"`
}

// Returns ID and Range keys
func GetSearchKeyByID(username, bookmarkId string) (string, string) {
	return fmt.Sprintf("USERNAME_%s", username), fmt.Sprintf("BOOKMARK_%s", bookmarkId)
}

// Returns ID and Range keys
func GetSearchKeyByName(username, name string) (string, string) {
	return fmt.Sprintf("USERNAME_%s", username), fmt.Sprintf("NAME_%s", name)
}

// Returns ID and Range keys
func GetSearchKeyByTag(username, tag string) (string, string) {
	return fmt.Sprintf("USERNAME_%s", username), fmt.Sprintf("TAG_%s", tag)
}

func (b *Bookmark) GetEntity() Bookmark {
	return Bookmark{
		Username:  fmt.Sprintf("USERNAME_%s", b.Username),
		ID:        fmt.Sprintf("BOOKMARK_%s", b.ID),
		Name:      b.Name,
		Url:       b.Url,
		Tags:      b.Tags,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

func (b *Bookmark) GetUsername() string {
	if strings.HasPrefix(b.Username, "USERNAME_") {
		return strings.Split(b.Username, "_")[1]
	}

	return b.Username
}

func (b *Bookmark) GetBookmarkId() string {
	if strings.HasPrefix(b.ID, "BOOKMARK_") {
		return strings.Split(b.ID, "_")[1]
	}

	return b.ID
}

func (b *Bookmark) GetSearchByName() BookmarkSearchByName {
	return BookmarkSearchByName{
		Username: fmt.Sprintf("USERNAME_%s", b.Username),
		Name:     fmt.Sprintf("NAME_%s_%s", b.Name, b.ID),
	}
}

func (b *Bookmark) GetSearchByTag() []BookmarkSearchByTag {
	return funk.Map(b.Tags, func(tag string) BookmarkSearchByTag {
		return BookmarkSearchByTag{
			Username: fmt.Sprintf("USERNAME_%s", b.Username),
			Tag:      fmt.Sprintf("TAG_%s_%s", tag, b.ID),
		}
	}).([]BookmarkSearchByTag)
}

// SearchByName
type BookmarkSearchByName struct {
	Username string `json:"username" dynamo:"id"`
	Name     string `json:"name" dynamo:"range"`
}

// Returns bookmarkId. Range key format NAME_{NAME}_{BOOKMARK_ID}
func (b *BookmarkSearchByName) GetBookmarkId() string {
	return strings.Split(b.Name, "_")[2]
}

func NewBookmarkSearchByTag(username, bookmarkId, tag string) BookmarkSearchByTag {
	return BookmarkSearchByTag{
		Username: fmt.Sprintf("USERNAME_%s", username),
		Tag:      fmt.Sprintf("TAG_%s_%s", tag, bookmarkId),
	}
}

// SearchByTag
type BookmarkSearchByTag struct {
	Username string `json:"username" dynamo:"id"`
	Tag      string `json:"tag" dynamo:"range"`
}

// Returns bookmarkId. Range key format TAG_{NAME}_{BOOKMARK_ID}
func (b *BookmarkSearchByTag) GetBookmarkId() string {
	return strings.Split(b.Tag, "_")[2]
}
