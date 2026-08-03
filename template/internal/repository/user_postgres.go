package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"{{ module_name }}/internal/apperror"
	"{{ module_name }}/internal/domain"
)

type postgresUserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	const q = `
		INSERT INTO users (id, name, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, q, user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	return mapPgError(err)
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const q = `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	return r.scanOne(r.db.QueryRow(ctx, q, id))
}

func (r *postgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	const q = `
		SELECT id, name, email, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	return r.scanOne(r.db.QueryRow(ctx, q, email))
}

func (r *postgresUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, int, error) {
	const countQ = `SELECT count(*) FROM users`
	var total int
	if err := r.db.QueryRow(ctx, countQ).Scan(&total); err != nil {
		return nil, 0, mapPgError(err)
	}

	const q = `
		SELECT id, name, email, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, 0, mapPgError(err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, mapPgError(err)
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, mapPgError(err)
	}

	return users, total, nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	const q = `
		UPDATE users
		SET name = $2, email = $3, updated_at = $4
		WHERE id = $1
	`
	tag, err := r.db.Exec(ctx, q, user.ID, user.Name, user.Email, user.UpdatedAt)
	if err != nil {
		return mapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return apperror.NotFound("user not found")
	}
	return nil
}

func (r *postgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const q = `DELETE FROM users WHERE id = $1`
	tag, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return mapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return apperror.NotFound("user not found")
	}
	return nil
}

func (r *postgresUserRepository) scanOne(row interface {
	Scan(dest ...any) error
}) (*domain.User, error) {
	var u domain.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, mapPgError(err)
	}
	return &u, nil
}
