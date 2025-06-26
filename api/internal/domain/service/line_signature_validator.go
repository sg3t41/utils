package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

// SignatureValidator 署名検証のインターフェース
type SignatureValidator interface {
	Validate(signature string, body []byte) bool
}

// signatureValidator 署名検証の実装
type signatureValidator struct {
	channelSecret string
}

// NewSignatureValidator 署名検証器の生成
func NewSignatureValidator(channelSecret string) SignatureValidator {
	return &signatureValidator{
		channelSecret: channelSecret,
	}
}

// Validate 署名を検証
func (v *signatureValidator) Validate(signature string, body []byte) bool {
	// プレフィックスを除去
	signature = strings.TrimPrefix(signature, SignaturePrefix)
	
	// 期待される署名を計算
	expectedSignature := v.calculateSignature(body)
	
	// 比較
	return signature == expectedSignature
}

// calculateSignature 署名を計算
func (v *signatureValidator) calculateSignature(body []byte) string {
	h := hmac.New(sha256.New, []byte(v.channelSecret))
	h.Write(body)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}