package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skroda/backend/internal/models"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrDuplicatePhone = errors.New("phone number already registered")
	ErrDuplicateEmail = errors.New("email already registered")
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, phone, password_hash, first_name, last_name, role, status, city, region)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at`

	user.ID = uuid.New()
	if user.Status == "" {
		user.Status = models.StatusPendingVerification
	}

	err := r.pool.QueryRow(ctx, query,
		user.ID, user.Email, user.Phone, user.PasswordHash,
		user.FirstName, user.LastName, user.Role, user.Status,
		user.City, user.Region,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "duplicate key") {
			if strings.Contains(msg, "users_phone_key") {
				return ErrDuplicatePhone
			}
			if strings.Contains(msg, "users_email_key") {
				return ErrDuplicateEmail
			}
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, phone, phone_verified, password_hash, first_name, last_name,
		       role, status, city, region, avatar_url, trust_score, total_transactions,
		       created_at, updated_at
		FROM users WHERE id = $1`

	user := &models.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Phone, &user.PhoneVerified,
		&user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.Status, &user.City, &user.Region,
		&user.AvatarURL, &user.TrustScore, &user.TotalTransactions,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := `
		SELECT id, email, phone, phone_verified, password_hash, first_name, last_name,
		       role, status, city, region, avatar_url, trust_score, total_transactions,
		       created_at, updated_at
		FROM users WHERE phone = $1`

	user := &models.User{}
	err := r.pool.QueryRow(ctx, query, phone).Scan(
		&user.ID, &user.Email, &user.Phone, &user.PhoneVerified,
		&user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.Status, &user.City, &user.Region,
		&user.AvatarURL, &user.TrustScore, &user.TotalTransactions,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, phone, phone_verified, password_hash, first_name, last_name,
		       role, status, city, region, avatar_url, trust_score, total_transactions,
		       created_at, updated_at
		FROM users WHERE email = $1`

	user := &models.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Phone, &user.PhoneVerified,
		&user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Role, &user.Status, &user.City, &user.Region,
		&user.AvatarURL, &user.TrustScore, &user.TotalTransactions,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *UserRepository) ListByRoleAndCity(ctx context.Context, role models.UserRole, city string) ([]models.User, error) {
	query := `
		SELECT id, email, phone, phone_verified, password_hash, first_name, last_name,
		       role, status, city, region, avatar_url, trust_score, total_transactions,
		       created_at, updated_at
		FROM users WHERE role = $1`
	args := []interface{}{role}

	if city != "" {
		query += " AND city = $2"
		args = append(args, city)
	}
	query += " ORDER BY trust_score DESC"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID, &u.Email, &u.Phone, &u.PhoneVerified,
			&u.PasswordHash, &u.FirstName, &u.LastName,
			&u.Role, &u.Status, &u.City, &u.Region,
			&u.AvatarURL, &u.TrustScore, &u.TotalTransactions,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) UpdateTrustScore(ctx context.Context, userID uuid.UUID, score float64, totalTxns int) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET trust_score = $1, total_transactions = $2, updated_at = NOW() WHERE id = $3`,
		score, totalTxns, userID)
	return err
}
