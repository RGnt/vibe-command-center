package repository

import (
	"database/sql"
	"errors"
	"todo-backend/internal/models"
)

// LibraryRepository ...
type LibraryRepository interface {
	CreateDocument(doc models.LibraryDocument) (models.LibraryDocument, error)
	GetDocumentsByUserID(userID int) ([]models.LibraryDocument, error)
	GetDocumentByID(id int, userID int) (models.LibraryDocument, error)
	DeleteDocument(id int, userID int) error
}

// PostgresLibraryRepository ...
type PostgresLibraryRepository struct {
	db *sql.DB
}

// NewLibraryRepository ...
func NewLibraryRepository(db *sql.DB) LibraryRepository {
	return &PostgresLibraryRepository{db: db}
}

// CreateDocument ...
func (r *PostgresLibraryRepository) CreateDocument(doc models.LibraryDocument) (models.LibraryDocument, error) {
	query := `
		INSERT INTO library_documents (user_id, original_name, filename, filepath, size, mime_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(query, doc.UserID, doc.OriginalName, doc.Filename, doc.Filepath, doc.Size, doc.MimeType).
		Scan(&doc.ID, &doc.CreatedAt)
	if err != nil {
		return doc, err
	}
	return doc, nil
}

// GetDocumentsByUserID ...
func (r *PostgresLibraryRepository) GetDocumentsByUserID(userID int) ([]models.LibraryDocument, error) {
	query := `
		SELECT id, user_id, original_name, filename, filepath, size, mime_type, created_at
		FROM library_documents
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var docs []models.LibraryDocument
	for rows.Next() {
		var doc models.LibraryDocument
		if err := rows.Scan(&doc.ID, &doc.UserID, &doc.OriginalName, &doc.Filename, &doc.Filepath, &doc.Size, &doc.MimeType, &doc.CreatedAt); err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

// GetDocumentByID ...
func (r *PostgresLibraryRepository) GetDocumentByID(id int, userID int) (models.LibraryDocument, error) {
	query := `
		SELECT id, user_id, original_name, filename, filepath, size, mime_type, created_at
		FROM library_documents
		WHERE id = $1 AND user_id = $2
	`
	var doc models.LibraryDocument
	err := r.db.QueryRow(query, id, userID).
		Scan(&doc.ID, &doc.UserID, &doc.OriginalName, &doc.Filename, &doc.Filepath, &doc.Size, &doc.MimeType, &doc.CreatedAt)
	if err == sql.ErrNoRows {
		return doc, errors.New("document not found")
	}
	return doc, err
}

// DeleteDocument ...
func (r *PostgresLibraryRepository) DeleteDocument(id int, userID int) error {
	query := `DELETE FROM library_documents WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("document not found or unauthorized")
	}
	return nil
}
