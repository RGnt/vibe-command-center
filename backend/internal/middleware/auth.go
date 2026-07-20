package middleware

import (
	"context"
	"fmt"
	"net/http"

	"todo-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the JWT payload validated by this middleware.
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// UserContextKey is the context key type for the authenticated user ID.
type UserContextKey string

// UserIDKey is the context key used to store and retrieve the authenticated user ID.
const UserIDKey UserContextKey = "user_id"

// AuthMiddlewareProvider holds dependencies for the auth middleware.
type AuthMiddlewareProvider struct {
	jwtKey    []byte
	blocklist repository.TokenBlocklistRepository
}

// NewAuthMiddlewareProvider creates an AuthMiddlewareProvider.
// blocklist is used to reject tokens that have been explicitly revoked (e.g., after logout).
func NewAuthMiddlewareProvider(jwtKey []byte, blocklist repository.TokenBlocklistRepository) *AuthMiddlewareProvider {
	return &AuthMiddlewareProvider{jwtKey: jwtKey, blocklist: blocklist}
}

// Middleware validates the JWT token, asserts the signing algorithm, checks the
// revocation blocklist, and injects the authenticated user ID into the request context.
func (m *AuthMiddlewareProvider) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Missing token cookie", http.StatusUnauthorized)
			return
		}
		tokenString := cookie.Value

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// AUTH-02: Assert the signing algorithm is exactly HMAC-SHA256.
			// Accepting any algorithm would allow algorithm-confusion attacks (e.g., alg:none).
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return m.jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// AUTH-01: Check that this token's jti has not been explicitly revoked
		// (e.g., the user has logged out).
		if claims.ID != "" {
			revoked, err := m.blocklist.IsRevoked(claims.ID)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			if revoked {
				http.Error(w, "Token has been revoked", http.StatusUnauthorized)
				return
			}
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts the authenticated user ID from the request context.
func GetUserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(UserIDKey).(int)
	return id, ok
}
