package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skroda/backend/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (id, full_name, email, phone, password_hash, role, city, country, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		user.ID, user.FullName, user.Email, user.Phone, user.PasswordHash,
		user.Role, user.City, user.Country, user.IsVerified, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, full_name, email, phone, password_hash, role, city, country, is_verified, created_at, updated_at
		FROM users WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&user.ID, &user.FullName, &user.Email, &user.Phone, &user.PasswordHash,
			&user.Role, &user.City, &user.Country, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, full_name, email, phone, password_hash, role, city, country, is_verified, created_at, updated_at
		FROM users WHERE email = $1 AND deleted_at IS NULL`, email).
		Scan(&user.ID, &user.FullName, &user.Email, &user.Phone, &user.PasswordHash,
			&user.Role, &user.City, &user.Country, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) ListByRoleAndCity(ctx context.Context, role models.UserRole, city string) ([]models.User, error) {
	query := `SELECT id, full_name, email, phone, role, city, country, is_verified, created_at, updated_at
		FROM users WHERE role = $1 AND deleted_at IS NULL`
	args := []interface{}{role}

	if city != "" {
		query += " AND city = $2"
		args = append(args, city)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.FullName, &u.Email, &u.Phone, &u.Role, &u.City, &u.Country, &u.IsVerified, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func parseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

var ErrNotFound = errors.New("record not found")
