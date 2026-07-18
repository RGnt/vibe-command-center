package handlers

import (
	"encoding/json"
	"net/http"

	"todo-backend/database"
	"todo-backend/middleware"
	"todo-backend/models"
)

func GetUserSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var settings models.UserSettings
	err := database.DB.QueryRow(`
		SELECT user_id, theme, default_project_id 
		FROM user_settings WHERE user_id = $1
	`, userID).Scan(&settings.UserID, &settings.Theme, &settings.DefaultProjectID)
	
	if err != nil {
		// If not found, return default
		settings = models.UserSettings{
			UserID: userID,
			Theme:  "system",
		}
	}

	json.NewEncoder(w).Encode(settings)
}

func UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var settings models.UserSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := database.DB.Exec(`
		INSERT INTO user_settings (user_id, theme, default_project_id) 
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET 
			theme = EXCLUDED.theme,
			default_project_id = EXCLUDED.default_project_id
	`, userID, settings.Theme, settings.DefaultProjectID)

	if err != nil {
		http.Error(w, "Error updating settings", http.StatusInternalServerError)
		return
	}

	settings.UserID = userID
	json.NewEncoder(w).Encode(settings)
}
