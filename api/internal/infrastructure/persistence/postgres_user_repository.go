package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sg3t41/api/internal/domain/entity"
	"github.com/sg3t41/api/internal/domain/repository"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, name, email, line_user_id, profile_image, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.LineUserID, user.ProfileImage, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user entity.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user entity.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) FindByLineUserID(ctx context.Context, lineUserID string) (*entity.User, error) {
	query := `
		SELECT id, name, email, line_user_id, profile_image, created_at, updated_at
		FROM users
		WHERE line_user_id = $1
	`
	var user entity.User
	err := r.db.QueryRowContext(ctx, query, lineUserID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.LineUserID,
		&user.ProfileImage,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) FindAll(ctx context.Context) ([]*entity.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, rows.Err()
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, updated_at = $4
		WHERE id = $1
	`
	user.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email, user.UpdatedAt)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	query := `
		SELECT id, name, email, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, rows.Err()
}

func (r *PostgresUserRepository) FindWithOffsetPagination(ctx context.Context, limit, offset int, filter repository.PaginationFilter, sort repository.SortOption) (*repository.PaginationResult, error) {
	baseQuery := "SELECT id, name, email, created_at, updated_at FROM users"
	countQuery := "SELECT COUNT(*) FROM users"
	
	whereClause, args := r.buildWhereClause(filter)
	orderClause := r.buildOrderClause(sort)
	
	if whereClause != "" {
		baseQuery += " WHERE " + whereClause
		countQuery += " WHERE " + whereClause
	}
	
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}
	
	baseQuery += orderClause + " LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)
	
	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	
	return &repository.PaginationResult{
		Users: users,
		Total: total,
	}, rows.Err()
}

func (r *PostgresUserRepository) FindWithCursorPagination(ctx context.Context, limit int, cursor string, filter repository.PaginationFilter, sort repository.SortOption) ([]*entity.User, error) {
	baseQuery := "SELECT id, name, email, created_at, updated_at FROM users"
	
	whereClause, args := r.buildWhereClause(filter)
	
	if cursor != "" {
		if whereClause != "" {
			whereClause += " AND created_at < $" + fmt.Sprintf("%d", len(args)+1)
		} else {
			whereClause = "created_at < $" + fmt.Sprintf("%d", len(args)+1)
		}
		args = append(args, cursor)
	}
	
	if whereClause != "" {
		baseQuery += " WHERE " + whereClause
	}
	
	orderClause := r.buildOrderClause(sort)
	baseQuery += orderClause + " LIMIT $" + fmt.Sprintf("%d", len(args)+1)
	args = append(args, limit+1)
	
	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	
	return users, rows.Err()
}

func (r *PostgresUserRepository) buildWhereClause(filter repository.PaginationFilter) (string, []interface{}) {
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

func (r *PostgresUserRepository) buildOrderClause(sort repository.SortOption) string {
	if sort.Field == "" {
		return " ORDER BY created_at DESC"
	}
	
	field := sort.Field
	if field != "id" && field != "name" && field != "email" && field != "created_at" && field != "updated_at" {
		field = "created_at"
	}
	
	order := sort.Order
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	
	return fmt.Sprintf(" ORDER BY %s %s", field, strings.ToUpper(order))
}

func (r *PostgresUserRepository) SoftDelete(ctx context.Context, id string) error {
	query := `
		UPDATE users
		SET deleted_at = $2, updated_at = $2
		WHERE id = $1 AND deleted_at IS NULL
	`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, id, now)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostgresUserRepository) HardDelete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}