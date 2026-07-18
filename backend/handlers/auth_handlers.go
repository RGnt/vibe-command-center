package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"todo-backend/database"
	"todo-backend/middleware"
	"todo-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var input AuthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	var user models.User
	err = database.DB.QueryRow(`
		INSERT INTO users (email, password_hash) 
		VALUES ($1, $2) RETURNING id, email, created_at
	`, input.Email, string(hash)).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		http.Error(w, "Email might already exist", http.StatusConflict)
		return
	}

	// Create default user settings
	_, _ = database.DB.Exec(`INSERT INTO user_settings (user_id, theme) VALUES ($1, 'system')`, user.ID)

	// Create a default workflow
	var workflowID int
	err = database.DB.QueryRow(`INSERT INTO workflows (user_id, name) VALUES ($1, 'Default Workflow') RETURNING id`, user.ID).Scan(&workflowID)
	if err == nil {
		stages := []string{"To Do", "In Progress", "Review", "Done"}
		for i, stageName := range stages {
			_, _ = database.DB.Exec(`INSERT INTO workflow_stages (workflow_id, name, "order") VALUES ($1, $2, $3)`, workflowID, stageName, i+1)
		}
		
		// Create a default project
		var projectID int
		err = database.DB.QueryRow(`INSERT INTO projects (user_id, name, description, workflow_id) VALUES ($1, 'General Project', 'Default project', $2) RETURNING id`, user.ID, workflowID).Scan(&projectID)
		if err == nil {
			// update user setting default project
			_, _ = database.DB.Exec(`UPDATE user_settings SET default_project_id = $1 WHERE user_id = $2`, projectID, user.ID)

			// Add a default wiki page
			_, _ = database.DB.Exec(`INSERT INTO wiki_pages (user_id, project_id, category, title, slug, content) VALUES ($1, $2, 'General', 'Welcome to the Wiki', 'welcome', '# Welcome\n\nThis is your private wiki.')`, user.ID, projectID)
		}
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{Token: token, User: user})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input AuthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user models.User
	var hash string
	err := database.DB.QueryRow(`SELECT id, email, password_hash, created_at FROM users WHERE email = $1`, input.Email).Scan(&user.ID, &user.Email, &hash, &user.CreatedAt)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AuthResponse{Token: token, User: user})
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	err := database.DB.QueryRow(`SELECT id, email, created_at FROM users WHERE id = $1`, userID).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func generateJWT(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * 7 * time.Hour) // 1 week
	claims := &middleware.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(middleware.JwtKey)
}
