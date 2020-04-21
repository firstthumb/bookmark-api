package entity

type Tag struct {
	Tag        string `json:"tag" dynamo:"tag"`
	BookmarkID string `json:"bookmark_id" dynamo:"bookmark_id"`
}
