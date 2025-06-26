package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/sg3t41/api/internal/domain/entity"
)

// LinkRepository リンクリポジトリインターフェース
type LinkRepository interface {
	// GetAll 全リンクを取得（表示順でソート）
	GetAll(ctx context.Context, userID uuid.UUID) ([]*entity.Link, error)
	
	// GetByID IDでリンクを取得
	GetByID(ctx context.Context, id int) (*entity.Link, error)
	
	// Create リンクを作成
	Create(ctx context.Context, link *entity.Link) error
	
	// Update リンクを更新
	Update(ctx context.Context, link *entity.Link) error
	
	// Delete リンクを削除
	Delete(ctx context.Context, id int) error
	
	// GetActiveLinks アクティブなリンクのみを取得
	GetActiveLinks(ctx context.Context, userID uuid.UUID) ([]*entity.Link, error)
	
	// UpdateOrder 表示順序を更新
	UpdateOrder(ctx context.Context, linkID int, newOrder int) error
}