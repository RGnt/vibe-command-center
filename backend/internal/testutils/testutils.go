package testutils

import (
	"context"
	"net/http"
	"time"

	"todo-backend/internal/database"
	mymiddleware "todo-backend/internal/middleware"

	"github.com/golang-jwt/jwt/v5"
)

// JWTSecret used in testing
var TestJWTSecret = []byte("super_secret_test_key")

// SetupTestDB initializes the test database and truncates users table
// which cascades and clears all other tables for a fresh state.
func SetupTestDB() {
	database.InitTestDB()
	ClearDB()
}

// ClearDB wipes all data from the database
func ClearDB() {
	// TRUNCATE CASCADE will clear all tables dependent on users
	_, err := database.DB.Exec(`TRUNCATE TABLE users RESTART IDENTITY CASCADE`)
	if err != nil {
		panic("Failed to truncate database: " + err.Error())
	}
}

// GenerateTestJWT generates a valid JWT token for a given userID
func GenerateTestJWT(userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(TestJWTSecret)
	return tokenString
}

// AuthContext middleware injects a specific user ID into the context
// This is useful for testing handlers directly without hitting the HTTP middleware
func AuthContext(userID int, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), mymiddleware.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
