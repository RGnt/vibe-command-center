package service

import (
	"errors"
	"time"

	"todo-backend/models"
	"todo-backend/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

// Claims matches the structure in middleware
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Register(email, password string) (models.User, string, error)
	Login(email, password string) (models.User, string, error)
	GetUserByID(id int) (models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwtKey   []byte
}

func NewAuthService(userRepo repository.UserRepository, jwtKey []byte) AuthService {
	return &authService{
		userRepo: userRepo,
		jwtKey:   jwtKey,
	}
}

func (s *authService) Register(email, password string) (models.User, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, "", err
	}

	user, err := s.userRepo.CreateUser(email, string(hash))
	if err != nil {
		return models.User{}, "", err
	}

	token, err := s.generateJWT(user.ID)
	if err != nil {
		return user, "", err
	}

	return user, token, nil
}

func (s *authService) Login(email, password string) (models.User, string, error) {
	user, hash, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return models.User{}, "", ErrInvalidCredentials
		}
		return models.User{}, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return models.User{}, "", ErrInvalidCredentials
	}

	token, err := s.generateJWT(user.ID)
	if err != nil {
		return user, "", err
	}

	return user, token, nil
}

func (s *authService) GetUserByID(id int) (models.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *authService) generateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * 7 * time.Hour) // 1 week
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtKey)
}
