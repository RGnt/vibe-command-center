package repository

import (
	"database/sql"
	"time"
)

// TokenBlocklistRepository provides storage for revoked JWT IDs.
// Once a jti is present in the store, the associated token is rejected by middleware.
type TokenBlocklistRepository interface {
	// RevokeToken inserts the token's jti into the blocklist until expiresAt.
	RevokeToken(jti string, expiresAt time.Time) error
	// IsRevoked returns true if the given jti is present in the blocklist.
	IsRevoked(jti string) (bool, error)
	// DeleteExpired removes blocklist entries whose tokens have already expired.
	// Intended to be called periodically to prevent unbounded table growth.
	DeleteExpired() error
}

type postgresTokenBlocklistRepository struct {
	db *sql.DB
}

// NewTokenBlocklistRepository creates a new Postgres-backed TokenBlocklistRepository.
func NewTokenBlocklistRepository(db *sql.DB) TokenBlocklistRepository {
	return &postgresTokenBlocklistRepository{db: db}
}

// RevokeToken inserts a jti into the revoked_tokens table.
func (r *postgresTokenBlocklistRepository) RevokeToken(jti string, expiresAt time.Time) error {
	_, err := r.db.Exec(
		`INSERT INTO revoked_tokens (jti, expires_at) VALUES ($1, $2) ON CONFLICT (jti) DO NOTHING`,
		jti, expiresAt,
	)
	return err
}

// IsRevoked checks if the jti exists in the revoked_tokens table.
func (r *postgresTokenBlocklistRepository) IsRevoked(jti string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM revoked_tokens WHERE jti = $1 AND expires_at > NOW())`,
		jti,
	).Scan(&exists)
	return exists, err
}

// DeleteExpired removes entries whose token expiry has already passed.
func (r *postgresTokenBlocklistRepository) DeleteExpired() error {
	_, err := r.db.Exec(`DELETE FROM revoked_tokens WHERE expires_at <= NOW()`)
	return err
}
