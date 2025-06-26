package service

// LINE API関連の定数
const (
	// LineAPIBaseURL LINE Messaging APIのベースURL
	LineAPIBaseURL = "https://api.line.me/v2/bot"
	
	// LineReplyEndpoint 返信エンドポイント
	LineReplyEndpoint = "/message/reply"
	
	// SignatureHeader 署名検証用のヘッダー名
	SignatureHeader = "X-Line-Signature"
	
	// SignaturePrefix 署名のプレフィックス
	SignaturePrefix = "sha256="
	
	// ContentTypeJSON JSONコンテンツタイプ
	ContentTypeJSON = "application/json"
)

// イベントタイプ定数
const (
	EventTypeMessage = "message"
)

// メッセージタイプ定数
const (
	MessageTypeText = "text"
)