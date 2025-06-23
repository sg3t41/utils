package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sg3t41/api/internal/interfaces/dto"
)

type LineService struct {
	config dto.LineOAuthConfig
	client *http.Client
}

func NewLineService(config dto.LineOAuthConfig) *LineService {
	return &LineService{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// LINE認証用URL生成
func (s *LineService) GetAuthURL(state string) string {
	baseURL := "https://access.line.me/oauth2/v2.1/authorize"
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", s.config.ClientID)
	params.Add("redirect_uri", s.config.RedirectURL)
	params.Add("state", state)
	params.Add("scope", strings.Join(s.config.Scope, " "))

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// アクセストークン取得
func (s *LineService) GetAccessToken(code string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.config.RedirectURL)
	data.Set("client_id", s.config.ClientID)
	data.Set("client_secret", s.config.ClientSecret)

	req, err := http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get access token: %s", string(body))
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Scope        string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

// ユーザープロフィール取得
func (s *LineService) GetProfile(accessToken string) (*dto.LineProfile, error) {
	req, err := http.NewRequest("GET", "https://api.line.me/v2/profile", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get profile: %s", string(body))
	}

	var profile dto.LineProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}