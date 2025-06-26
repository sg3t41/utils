package persistence

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sg3t41/api/internal/domain/repository"
)

// NewSqlxDB sql.DBからsqlx.DBを作成
func NewSqlxDB(db *sql.DB) *sqlx.DB {
	return sqlx.NewDb(db, "postgres")
}

// 共通のクエリビルダー機能

// buildWhereClause 共通のWHERE句構築
func buildWhereClause(filter repository.PaginationFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR email ILIKE $%d)", argIndex, argIndex+1))
		searchPattern := "%" + filter.Search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	if filter.Status != "" {
		switch filter.Status {
		case "active":
			conditions = append(conditions, "deleted_at IS NULL")
		case "deleted":
			conditions = append(conditions, "deleted_at IS NOT NULL")
		}
	}

	if filter.CreatedFrom != "" {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, filter.CreatedFrom)
		argIndex++
	}

	if filter.CreatedTo != "" {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, filter.CreatedTo)
		argIndex++
	}

	return strings.Join(conditions, " AND "), args
}

// buildOrderClause 共通のORDER BY句構築
func buildOrderClause(sort repository.SortOption) string {
	if sort.Field == "" {
		return " ORDER BY created_at DESC"
	}

	field := sort.Field
	// セキュリティのため、許可されたフィールドのみ使用
	allowedFields := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"created_at": true,
		"updated_at": true,
	}
	if !allowedFields[field] {
		field = "created_at"
	}

	order := sort.Order
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	return fmt.Sprintf(" ORDER BY %s %s", field, strings.ToUpper(order))
}

// buildPaginationClause ページネーション句構築
func buildPaginationClause(limit, offset int, currentArgIndex int) (string, []interface{}) {
	clause := fmt.Sprintf(" LIMIT $%d OFFSET $%d", currentArgIndex, currentArgIndex+1)
	args := []interface{}{limit, offset}
	return clause, args
}