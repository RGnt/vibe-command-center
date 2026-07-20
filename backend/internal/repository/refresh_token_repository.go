package repository

import (
	"database/sql"
	"errors"
	"time"
)

// ErrInvalidRefreshToken is returned when a refresh token hash is not found or expired.
var ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")

// RefreshToken represents a stored refresh token hash and its expiry.
type RefreshToken struct {
	ID        int
	UserID    int
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// RefreshTokenRepository defines the interface for refresh token storage.
type RefreshTokenRepository interface {
	Create(userID int, tokenHash string, expiresAt time.Time) error
	FindByHash(tokenHash string) (RefreshToken, error)
	DeleteByHash(tokenHash string) error
	DeleteExpired() error
}

type postgresRefreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository creates a new Postgres-backed RefreshTokenRepository.
func NewRefreshTokenRepository(db *sql.DB) RefreshTokenRepository {
	return &postgresRefreshTokenRepository{db: db}
}

// Create stores a new refresh token hash for a user.
func (r *postgresRefreshTokenRepository) Create(userID int, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt,
	)
	return err
}

// FindByHash retrieves a refresh token record by its hash.
func (r *postgresRefreshTokenRepository) FindByHash(tokenHash string) (RefreshToken, error) {
	var rt RefreshToken
	err := r.db.QueryRow(
		`SELECT id, user_id, token_hash, expires_at, created_at FROM refresh_tokens WHERE token_hash = $1`,
		tokenHash,
	).Scan(&rt.ID, &rt.UserID, &rt.TokenHash, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return rt, ErrInvalidRefreshToken
		}
		return rt, err
	}
	return rt, nil
}

// DeleteByHash removes a refresh token from the database.
func (r *postgresRefreshTokenRepository) DeleteByHash(tokenHash string) error {
	_, err := r.db.Exec(`DELETE FROM refresh_tokens WHERE token_hash = $1`, tokenHash)
	return err
}

// DeleteExpired removes entries whose token expiry has passed.
func (r *postgresRefreshTokenRepository) DeleteExpired() error {
	_, err := r.db.Exec(`DELETE FROM refresh_tokens WHERE expires_at <= NOW()`)
	return err
}
