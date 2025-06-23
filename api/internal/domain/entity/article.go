package entity

import (
	"encoding/json"
	"time"
)

// ArticleStatus represents the status of an article
type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "draft"
	ArticleStatusPublished ArticleStatus = "published"
	ArticleStatusArchived  ArticleStatus = "archived"
)

// Article represents a blog article
type Article struct {
	ID             string        `json:"id" db:"id"`
	Title          string        `json:"title" db:"title"`
	Content        string        `json:"content" db:"content"`
	Summary        string        `json:"summary" db:"summary"`
	Status         ArticleStatus `json:"status" db:"status"`
	AuthorID       string        `json:"author_id" db:"author_id"`
	Tags           []string      `json:"tags" db:"tags"`
	ArticleImage   *string       `json:"article_image" db:"article_image"`
	CreatedAt      time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" db:"updated_at"`
	PublishedAt    *time.Time    `json:"published_at" db:"published_at"`
}

// IsPublished returns true if the article is published
func (a *Article) IsPublished() bool {
	return a.Status == ArticleStatusPublished
}

// IsDraft returns true if the article is in draft status
func (a *Article) IsDraft() bool {
	return a.Status == ArticleStatusDraft
}

// IsArchived returns true if the article is archived
func (a *Article) IsArchived() bool {
	return a.Status == ArticleStatusArchived
}

// Publish sets the article status to published and sets published_at timestamp
func (a *Article) Publish() {
	a.Status = ArticleStatusPublished
	now := time.Now()
	a.PublishedAt = &now
	a.UpdatedAt = now
}

// Unpublish sets the article status to draft and clears published_at timestamp
func (a *Article) Unpublish() {
	a.Status = ArticleStatusDraft
	a.PublishedAt = nil
	a.UpdatedAt = time.Now()
}

// Archive sets the article status to archived
func (a *Article) Archive() {
	a.Status = ArticleStatusArchived
	a.UpdatedAt = time.Now()
}

// HasTag returns true if the article has the specified tag
func (a *Article) HasTag(tag string) bool {
	for _, t := range a.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// AddTag adds a tag to the article if it doesn't already exist
func (a *Article) AddTag(tag string) {
	if !a.HasTag(tag) {
		a.Tags = append(a.Tags, tag)
		a.UpdatedAt = time.Now()
	}
}

// RemoveTag removes a tag from the article
func (a *Article) RemoveTag(tag string) {
	for i, t := range a.Tags {
		if t == tag {
			a.Tags = append(a.Tags[:i], a.Tags[i+1:]...)
			a.UpdatedAt = time.Now()
			break
		}
	}
}

// MarshalTags converts tags slice to JSON for database storage
func (a *Article) MarshalTags() ([]byte, error) {
	return json.Marshal(a.Tags)
}

// UnmarshalTags converts JSON from database to tags slice
func (a *Article) UnmarshalTags(data []byte) error {
	return json.Unmarshal(data, &a.Tags)
}

// SetArticleImage sets the article image path
func (a *Article) SetArticleImage(imagePath string) {
	a.ArticleImage = &imagePath
	a.UpdatedAt = time.Now()
}

// ClearImage removes the image reference
func (a *Article) ClearImage() {
	a.ArticleImage = nil
	a.UpdatedAt = time.Now()
}

// HasImage returns true if the article has an image
func (a *Article) HasImage() bool {
	return a.ArticleImage != nil && *a.ArticleImage != ""
}