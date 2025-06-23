package dto

// LINE ログインリクエスト
type LineLoginRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state,omitempty"`
}

// LINE ログインレスポンス
type LineLoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// LINE プロフィール情報
type LineProfile struct {
	UserID      string `json:"userId"`
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

// LINE OAuth設定
type LineOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scope        []string
}