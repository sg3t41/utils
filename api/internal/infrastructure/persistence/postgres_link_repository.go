package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type postgresLinkRepository struct {
	db *sqlx.DB
}

// NewPostgresLinkRepository リンクリポジトリの新しいインスタンスを作成
func NewPostgresLinkRepository(db *sqlx.DB) repository.LinkRepository {
	return &postgresLinkRepository{db: db}
}

// GetAll 全リンクを取得（表示順でソート）
func (r *postgresLinkRepository) GetAll(ctx context.Context, userID uuid.UUID) ([]*entity.Link, error) {
	query := `
		SELECT id, title, url, description, platform, icon_name, 
		       background_color, text_color, order_index, is_active, 
		       user_id, created_at, updated_at
		FROM links 
		WHERE user_id = $1 
		ORDER BY order_index ASC, created_at ASC`
	
	var links []*entity.Link
	err := r.db.SelectContext(ctx, &links, query, userID)
	if err != nil {
		return nil, fmt.Errorf("リンク一覧の取得に失敗: %w", err)
	}
	
	return links, nil
}

// GetByID IDでリンクを取得
func (r *postgresLinkRepository) GetByID(ctx context.Context, id int) (*entity.Link, error) {
	query := `
		SELECT id, title, url, description, platform, icon_name, 
		       background_color, text_color, order_index, is_active, 
		       user_id, created_at, updated_at
		FROM links 
		WHERE id = $1`
	
	var link entity.Link
	err := r.db.GetContext(ctx, &link, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("リンクが見つかりません: id=%d", id)
		}
		return nil, fmt.Errorf("リンクの取得に失敗: %w", err)
	}
	
	return &link, nil
}

// Create リンクを作成
func (r *postgresLinkRepository) Create(ctx context.Context, link *entity.Link) error {
	query := `
		INSERT INTO links (title, url, description, platform, icon_name, 
		                  background_color, text_color, order_index, is_active, user_id)
		VALUES (:title, :url, :description, :platform, :icon_name, 
		        :background_color, :text_color, :order_index, :is_active, :user_id)
		RETURNING id, created_at, updated_at`
	
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("リンク作成のクエリ準備に失敗: %w", err)
	}
	defer stmt.Close()
	
	err = stmt.GetContext(ctx, link, link)
	if err != nil {
		return fmt.Errorf("リンクの作成に失敗: %w", err)
	}
	
	return nil
}

// Update リンクを更新
func (r *postgresLinkRepository) Update(ctx context.Context, link *entity.Link) error {
	query := `
		UPDATE links 
		SET title = :title, url = :url, description = :description, 
		    platform = :platform, icon_name = :icon_name, 
		    background_color = :background_color, text_color = :text_color, 
		    order_index = :order_index, is_active = :is_active,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = :id AND user_id = :user_id
		RETURNING updated_at`
	
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return fmt.Errorf("リンク更新のクエリ準備に失敗: %w", err)
	}
	defer stmt.Close()
	
	err = stmt.GetContext(ctx, &link.UpdatedAt, link)
	if err != nil {
		return fmt.Errorf("リンクの更新に失敗: %w", err)
	}
	
	return nil
}

// Delete リンクを削除
func (r *postgresLinkRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM links WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("リンクの削除に失敗: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("削除結果の確認に失敗: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("削除対象のリンクが見つかりません: id=%d", id)
	}
	
	return nil
}

// GetActiveLinks アクティブなリンクのみを取得
func (r *postgresLinkRepository) GetActiveLinks(ctx context.Context, userID uuid.UUID) ([]*entity.Link, error) {
	query := `
		SELECT id, title, url, description, platform, icon_name, 
		       background_color, text_color, order_index, is_active, 
		       user_id, created_at, updated_at
		FROM links 
		WHERE user_id = $1 AND is_active = true 
		ORDER BY order_index ASC, created_at ASC`
	
	var links []*entity.Link
	err := r.db.SelectContext(ctx, &links, query, userID)
	if err != nil {
		return nil, fmt.Errorf("アクティブなリンク一覧の取得に失敗: %w", err)
	}
	
	return links, nil
}

// UpdateOrder 表示順序を更新
func (r *postgresLinkRepository) UpdateOrder(ctx context.Context, linkID int, newOrder int) error {
	query := `
		UPDATE links 
		SET order_index = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2`
	
	_, err := r.db.ExecContext(ctx, query, newOrder, linkID)
	if err != nil {
		return fmt.Errorf("表示順序の更新に失敗: %w", err)
	}
	
	return nil
}