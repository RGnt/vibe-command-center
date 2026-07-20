package models

import "time"

// LibraryDocument ...
type LibraryDocument struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	OriginalName string    `json:"original_name"`
	Filename     string    `json:"filename"`
	Filepath     string    `json:"-"`
	Size         int64     `json:"size"`
	MimeType     string    `json:"mime_type"`
	CreatedAt    time.Time `json:"created_at"`
}
