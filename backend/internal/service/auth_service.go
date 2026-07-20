package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"todo-backend/internal/models"
	"todo-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidCredentials is returned when login credentials do not match.
var ErrInvalidCredentials = errors.New("invalid credentials")

// Claims is the JWT payload structure used when signing and parsing tokens.
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// AuthService defines the contract for user authentication operations.
type AuthService interface {
	Register(email, password string) (models.User, string, string, error)
	Login(email, password string) (models.User, string, string, error)
	GetUserByID(id int) (models.User, error)
	RevokeToken(jti string, expiresAt time.Time) error
	RevokeTokenByString(tokenString string) error
	Refresh(refreshToken string) (string, string, error)
	Logout(refreshToken string) error
}

type authService struct {
	userRepo   repository.UserRepository
	blocklist  repository.TokenBlocklistRepository
	refreshRepo repository.RefreshTokenRepository
	jwtKey     []byte
	bcryptCost int
	dummyBcryptHash []byte
}

// NewAuthService creates a new AuthService with the given dependencies.
func NewAuthService(userRepo repository.UserRepository, blocklist repository.TokenBlocklistRepository, refreshRepo repository.RefreshTokenRepository, jwtKey []byte, bcryptCost int) AuthService {
	// Pre-compute a dummy hash to equalize timing (AUTH-06)
	dummyHash, _ := bcrypt.GenerateFromPassword([]byte("dummy_password_for_timing"), bcryptCost)

	return &authService{
		userRepo:        userRepo,
		blocklist:       blocklist,
		refreshRepo:     refreshRepo,
		jwtKey:          jwtKey,
		bcryptCost:      bcryptCost,
		dummyBcryptHash: dummyHash,
	}
}

// Register creates a new user account and returns the user, access token, and refresh token.
func (s *authService) Register(email, password string) (models.User, string, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return models.User{}, "", "", err
	}

	user, err := s.userRepo.CreateUser(email, string(hash))
	if err != nil {
		// AUTH-06: Equalize timing if email already exists
		if errors.Is(err, repository.ErrEmailExists) {
			_ = bcrypt.CompareHashAndPassword(s.dummyBcryptHash, []byte(password))
		}
		return models.User{}, "", "", err
	}

	accessToken, refreshToken, err := s.generateTokens(user.ID)
	return user, accessToken, refreshToken, err
}

// Login validates credentials and returns the user, access token, and refresh token.
func (s *authService) Login(email, password string) (models.User, string, string, error) {
	// AUTH-11: Normalize email for login consistency
	email = strings.TrimSpace(strings.ToLower(email))

	user, hash, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return models.User{}, "", "", ErrInvalidCredentials
		}
		return models.User{}, "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return models.User{}, "", "", ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.generateTokens(user.ID)
	return user, accessToken, refreshToken, err
}

// GetUserByID fetches a user record by its primary key.
func (s *authService) GetUserByID(id int) (models.User, error) {
	return s.userRepo.GetUserByID(id)
}

// RevokeToken writes the jti to the server-side blocklist, invalidating the token.
func (s *authService) RevokeToken(jti string, expiresAt time.Time) error {
	return s.blocklist.RevokeToken(jti, expiresAt)
}

// RevokeTokenByString parses a JWT and revokes its jti.
func (s *authService) RevokeTokenByString(tokenString string) error {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	
	if err != nil || !token.Valid || claims.ID == "" {
		return nil
	}
	
	return s.RevokeToken(claims.ID, claims.ExpiresAt.Time)
}

// Refresh validates a refresh token, rotates it, and returns a new access/refresh pair.
func (s *authService) Refresh(refreshToken string) (string, string, error) {
	hashBytes := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hashBytes[:])

	rt, err := s.refreshRepo.FindByHash(tokenHash)
	if err != nil {
		return "", "", err // Could be ErrInvalidRefreshToken
	}

	if time.Now().After(rt.ExpiresAt) {
		_ = s.refreshRepo.DeleteByHash(tokenHash)
		return "", "", repository.ErrInvalidRefreshToken
	}

	// Delete old token (rotate)
	_ = s.refreshRepo.DeleteByHash(tokenHash)

	return s.generateTokens(rt.UserID)
}

// Logout revokes the refresh token.
func (s *authService) Logout(refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	hashBytes := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hashBytes[:])
	return s.refreshRepo.DeleteByHash(tokenHash)
}

// generateTokens creates a signed HS256 JWT access token (15m) and a secure refresh token (30d).
func (s *authService) generateTokens(userID int) (string, string, error) {
	// Access Token (15m)
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", "", err
	}

	// Refresh Token (30d)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	refreshToken := hex.EncodeToString(bytes)
	
	hashBytes := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hashBytes[:])

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	err = s.refreshRepo.Create(userID, tokenHash, expiresAt)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
