package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type PostgresArticleRepository struct {
	db *sql.DB
}

func NewPostgresArticleRepository(db *sql.DB) repository.ArticleRepository {
	return &PostgresArticleRepository{db: db}
}

func (r *PostgresArticleRepository) Create(ctx context.Context, article *entity.Article) error {
	tagsJSON, err := json.Marshal(article.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		INSERT INTO articles (id, title, content, summary, status, tags, article_image, created_at, updated_at, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = r.db.ExecContext(ctx, query,
		article.ID,
		article.Title,
		article.Content,
		article.Summary,
		article.Status,
		tagsJSON,
		article.ArticleImage,
		article.CreatedAt,
		article.UpdatedAt,
		article.PublishedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}

	return nil
}

func (r *PostgresArticleRepository) FindByID(ctx context.Context, id string) (*entity.Article, error) {
	query := `
		SELECT id, title, content, summary, status, tags, article_image, created_at, updated_at, published_at
		FROM articles
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	article := &entity.Article{}
	var tagsJSON []byte

	err := row.Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.Summary,
		&article.Status,
		&tagsJSON,
		&article.ArticleImage,
		&article.CreatedAt,
		&article.UpdatedAt,
		&article.PublishedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("article not found")
		}
		return nil, fmt.Errorf("failed to find article: %w", err)
	}

	if err := json.Unmarshal(tagsJSON, &article.Tags); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	return article, nil
}

func (r *PostgresArticleRepository) FindAll(ctx context.Context, filter repository.ArticleFilter) ([]*entity.Article, int, error) {
	// Build WHERE clause
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}


	if filter.Tag != nil {
		conditions = append(conditions, fmt.Sprintf("tags @> $%d", argIndex))
		tagJSON, _ := json.Marshal([]string{*filter.Tag})
		args = append(args, tagJSON)
		argIndex++
	}

	if filter.Search != nil {
		searchPattern := "%" + *filter.Search + "%"
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR content ILIKE $%d)", argIndex, argIndex+1))
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	if filter.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filter.DateFrom)
		argIndex++
	}

	if filter.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filter.DateTo)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM articles %s", whereClause)
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	// Build main query with pagination and sorting
	query := fmt.Sprintf(`
		SELECT id, title, content, summary, status, tags, article_image, created_at, updated_at, published_at
		FROM articles
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, filter.GetSort(), filter.GetOrder(), argIndex, argIndex+1)

	args = append(args, filter.GetLimit(), filter.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find articles: %w", err)
	}
	defer rows.Close()

	var articles []*entity.Article
	for rows.Next() {
		article := &entity.Article{}
		var tagsJSON []byte

		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.Summary,
			&article.Status,
			&tagsJSON,
			&article.ArticleImage,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.PublishedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan article: %w", err)
		}

		if err := json.Unmarshal(tagsJSON, &article.Tags); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal tags: %w", err)
		}

		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating articles: %w", err)
	}

	return articles, total, nil
}

func (r *PostgresArticleRepository) Update(ctx context.Context, article *entity.Article) error {
	tagsJSON, err := json.Marshal(article.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	query := `
		UPDATE articles
		SET title = $2, content = $3, summary = $4, status = $5, tags = $6, article_image = $7, updated_at = $8, published_at = $9
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		article.ID,
		article.Title,
		article.Content,
		article.Summary,
		article.Status,
		tagsJSON,
		article.ArticleImage,
		article.UpdatedAt,
		article.PublishedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update article: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

func (r *PostgresArticleRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM articles WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

func (r *PostgresArticleRepository) FindByStatus(ctx context.Context, status entity.ArticleStatus, limit, offset int) ([]*entity.Article, int, error) {
	filter := repository.ArticleFilter{
		Status: &status,
		Limit:  limit,
		Page:   (offset / limit) + 1,
	}
	return r.FindAll(ctx, filter)
}

func (r *PostgresArticleRepository) FindByTag(ctx context.Context, tag string, limit, offset int) ([]*entity.Article, int, error) {
	filter := repository.ArticleFilter{
		Tag:   &tag,
		Limit: limit,
		Page:  (offset / limit) + 1,
	}
	return r.FindAll(ctx, filter)
}

