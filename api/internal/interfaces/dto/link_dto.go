package dto

// CreateLinkRequest リンク作成リクエスト
type CreateLinkRequest struct {
	Title           string  `json:"title" binding:"required,max=100"`
	URL             string  `json:"url" binding:"required,url"`
	Description     *string `json:"description,omitempty"`
	Platform        string  `json:"platform" binding:"required,max=50"`
	IconName        *string `json:"icon_name,omitempty"`
	BackgroundColor *string `json:"background_color,omitempty"`
	TextColor       *string `json:"text_color,omitempty"`
	OrderIndex      *int    `json:"order_index,omitempty"`
	IsActive        *bool   `json:"is_active,omitempty"`
}

// UpdateLinkRequest リンク更新リクエスト
type UpdateLinkRequest struct {
	Title           *string `json:"title,omitempty" binding:"omitempty,max=100"`
	URL             *string `json:"url,omitempty" binding:"omitempty,url"`
	Description     *string `json:"description,omitempty"`
	Platform        *string `json:"platform,omitempty" binding:"omitempty,max=50"`
	IconName        *string `json:"icon_name,omitempty"`
	BackgroundColor *string `json:"background_color,omitempty"`
	TextColor       *string `json:"text_color,omitempty"`
	OrderIndex      *int    `json:"order_index,omitempty"`
	IsActive        *bool   `json:"is_active,omitempty"`
}

// LinkResponse リンク情報のレスポンス
type LinkResponse struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	URL             string  `json:"url"`
	Description     *string `json:"description"`
	Platform        string  `json:"platform"`
	IconName        *string `json:"icon_name"`
	BackgroundColor string  `json:"background_color"`
	TextColor       string  `json:"text_color"`
	OrderIndex      int     `json:"order_index"`
	IsActive        bool    `json:"is_active"`
	UserID          string  `json:"user_id"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// UpdateOrderRequest 表示順序更新リクエスト
type UpdateOrderRequest struct {
	OrderIndex int `json:"order_index" binding:"required,min=0"`
}

// LinkListResponse リンク一覧レスポンス
type LinkListResponse struct {
	Links []*LinkResponse `json:"links"`
	Total int             `json:"total"`
}