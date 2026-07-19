package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// UserContextKey is the key for the user id in the request context
type UserContextKey string

const UserIDKey UserContextKey = "user_id"

// AuthMiddlewareProvider holds dependencies for the auth middleware
type AuthMiddlewareProvider struct {
	jwtKey []byte
}

func NewAuthMiddlewareProvider(jwtKey []byte) *AuthMiddlewareProvider {
	return &AuthMiddlewareProvider{jwtKey: jwtKey}
}

// Middleware validates the JWT token and extracts the user ID
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
			return m.jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts the user ID from the request context
func GetUserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(UserIDKey).(int)
	return id, ok
}
