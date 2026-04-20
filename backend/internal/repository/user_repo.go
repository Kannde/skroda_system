package repository

import (
	"context"

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

const userSelectCols = `id, email, phone, phone_verified, password_hash,
	first_name, last_name, role, status, city, region, avatar_url,
	trust_score, total_transactions, created_at, updated_at`

func scanUser(row interface{ Scan(...interface{}) error }) (*models.User, error) {
	u := &models.User{}
	err := row.Scan(
		&u.ID, &u.Email, &u.Phone, &u.PhoneVerified, &u.PasswordHash,
		&u.FirstName, &u.LastName, &u.Role, &u.Status, &u.City, &u.Region, &u.AvatarURL,
		&u.TrustScore, &u.TotalTransactions, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users
			(id, email, phone, password_hash, first_name, last_name, role, status, city, region, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		user.ID, user.Email, user.Phone, user.PasswordHash,
		user.FirstName, user.LastName, user.Role, user.Status,
		user.City, user.Region, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+userSelectCols+` FROM users WHERE id = $1`, id)
	return scanUser(row)
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+userSelectCols+` FROM users WHERE phone = $1`, phone)
	return scanUser(row)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+userSelectCols+` FROM users WHERE email = $1`, email)
	return scanUser(row)
}

func (r *UserRepository) ListByRoleAndCity(ctx context.Context, role models.UserRole, city string) ([]models.User, error) {
	query := `SELECT ` + userSelectCols + ` FROM users WHERE role = $1`
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
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *u)
	}
	return users, nil
}
