package repository

import (
	"database/sql"
	"time"
	"todo-backend/internal/models"
)

// WikiRepository ...
type WikiRepository interface {
	GetAllByUserID(userID int) ([]models.WikiPage, error)
	GetAllByProjectID(userID int, projectID int) ([]models.WikiPage, error)
	GetBySlugAndUserID(slug string, userID int) (models.WikiPage, error)
	Create(wiki models.WikiPage) (models.WikiPage, error)
	Update(id int, wiki models.WikiPage) (models.WikiPage, error)
	Delete(id int, userID int) error
}

type wikiRepository struct {
	db *sql.DB
}

// NewWikiRepository ...
func NewWikiRepository(db *sql.DB) WikiRepository {
	return &wikiRepository{db: db}
}

// GetAllByUserID ...
func (r *wikiRepository) GetAllByUserID(userID int) ([]models.WikiPage, error) {
	rows, err := r.db.Query(`SELECT id, project_id, category, title, slug, content, created_at, updated_at FROM wiki_pages WHERE user_id = $1 ORDER BY category, title`, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var wikis []models.WikiPage
	for rows.Next() {
		var wiki models.WikiPage
		var projectID sql.NullInt64
		var createdAt, updatedAt time.Time

		err := rows.Scan(&wiki.ID, &projectID, &wiki.Category, &wiki.Title, &wiki.Slug, &wiki.Content, &createdAt, &updatedAt)
		if err != nil {
			continue
		}

		wiki.UserID = userID
		wiki.CreatedAt = createdAt
		wiki.UpdatedAt = updatedAt
		
		if projectID.Valid {
			pid := int(projectID.Int64)
			wiki.ProjectID = &pid
		}

		wikis = append(wikis, wiki)
	}

	return wikis, nil
}

// GetBySlugAndUserID ...
func (r *wikiRepository) GetBySlugAndUserID(slug string, userID int) (models.WikiPage, error) {
	var wiki models.WikiPage
	var projectID sql.NullInt64
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(`SELECT id, project_id, category, title, slug, content, created_at, updated_at FROM wiki_pages WHERE slug = $1 AND user_id = $2`, slug, userID).
		Scan(&wiki.ID, &projectID, &wiki.Category, &wiki.Title, &wiki.Slug, &wiki.Content, &createdAt, &updatedAt)
	if err != nil {
		return wiki, err
	}

	wiki.UserID = userID
	wiki.CreatedAt = createdAt
	wiki.UpdatedAt = updatedAt

	if projectID.Valid {
		pid := int(projectID.Int64)
		wiki.ProjectID = &pid
	}

	return wiki, nil
}

// Create ...
func (r *wikiRepository) Create(wiki models.WikiPage) (models.WikiPage, error) {
	var id int64
	query := `INSERT INTO wiki_pages (user_id, project_id, category, title, slug, content) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRow(query, wiki.UserID, wiki.ProjectID, wiki.Category, wiki.Title, wiki.Slug, wiki.Content).Scan(&id)
	if err != nil {
		return wiki, err
	}
	
	wiki.ID = int(id)
	wiki.CreatedAt = time.Now()
	wiki.UpdatedAt = time.Now()
	
	return wiki, nil
}

// Update ...
func (r *wikiRepository) Update(id int, wiki models.WikiPage) (models.WikiPage, error) {
	query := `UPDATE wiki_pages SET category = $1, title = $2, slug = $3, content = $4, project_id = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6 AND user_id = $7`
	_, err := r.db.Exec(query, wiki.Category, wiki.Title, wiki.Slug, wiki.Content, wiki.ProjectID, id, wiki.UserID)
	if err != nil {
		return wiki, err
	}
	
	wiki.ID = id
	return wiki, nil
}

// Delete ...
func (r *wikiRepository) Delete(id int, userID int) error {
	_, err := r.db.Exec("DELETE FROM wiki_pages WHERE id = $1 AND user_id = $2", id, userID)
	return err
}
