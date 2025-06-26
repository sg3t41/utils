package entity

import (
	"time"
	
	"github.com/google/uuid"
)

// Link リンクエンティティ
type Link struct {
	ID              int       `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	URL             string    `json:"url" db:"url"`
	Description     *string   `json:"description" db:"description"`
	Platform        string    `json:"platform" db:"platform"`
	IconName        *string   `json:"icon_name" db:"icon_name"`
	BackgroundColor string    `json:"background_color" db:"background_color"`
	TextColor       string    `json:"text_color" db:"text_color"`
	OrderIndex      int       `json:"order_index" db:"order_index"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	UserID          uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// LinkPlatform プラットフォーム定数
type LinkPlatform string

const (
	PlatformTwitter   LinkPlatform = "twitter"
	PlatformInstagram LinkPlatform = "instagram"
	PlatformGitHub    LinkPlatform = "github"
	PlatformLINE      LinkPlatform = "line"
	PlatformWebsite   LinkPlatform = "website"
	PlatformYouTube   LinkPlatform = "youtube"
	PlatformTikTok    LinkPlatform = "tiktok"
	PlatformLinkedIn  LinkPlatform = "linkedin"
)

// GetDefaultColors プラットフォームのデフォルト色を取得
func (p LinkPlatform) GetDefaultColors() (background, text string) {
	switch p {
	case PlatformTwitter:
		return "#1DA1F2", "#FFFFFF"
	case PlatformInstagram:
		return "#E4405F", "#FFFFFF"
	case PlatformGitHub:
		return "#333333", "#FFFFFF"
	case PlatformLINE:
		return "#00C300", "#FFFFFF"
	case PlatformYouTube:
		return "#FF0000", "#FFFFFF"
	case PlatformTikTok:
		return "#000000", "#FFFFFF"
	case PlatformLinkedIn:
		return "#0A66C2", "#FFFFFF"
	default:
		return "#6B7280", "#FFFFFF"
	}
}